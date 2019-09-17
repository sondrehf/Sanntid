package orderHandler

import (
	"fmt"
	"time"

	"../elevio"
	"../fsm"
	"../message"
	"../network/bcast"
	"../order"
	"../supervisor"
)

//OrderServer is the heart of our project. Transmits distributed orders and alive signals to message-server
// and decides what happens to the received messages from other elevators by the message type
func OrderServer(
	elevatorID int,
	drv_button chan elevio.ButtonEvent,
	receive_broadcast_message chan message.Message,
	order_finished chan order.Order,
	alive_map chan map[int]bool,
	send_broadcast_message chan message.Message,
	broadcast_message chan message.Message,
	new_order chan order.Order,
	receive_alive chan int,
	set_lights chan order.Light,
	lights_finished chan order.Light,
	engine_failed chan bool) {

	//Variables
	udpPort := 16512
	var currentFloor int
	var msgID int = 0
	elevOnFloor := make(map[int]int)
	elevTaken := make(map[int]bool)

	//Goroutines independent of init and needed in init
	go bcast.Transmitter(udpPort, broadcast_message)
	go bcast.Receiver(udpPort, receive_broadcast_message)
	go fsm.FsmFunc(new_order, order_finished, set_lights, lights_finished, engine_failed)
	go message.MessageServer(elevatorID, send_broadcast_message, broadcast_message)

	//Init
	database, newAliveMap := initDatabaseAndAliveMap(elevatorID, alive_map, receive_alive, receive_broadcast_message)

	//Goroutines dependent on init
	go supervisor.SetElevatorsAlive(elevatorID, receive_alive, alive_map)
	go elevio.PollButtons(drv_button)

	//timer to send I'm alive signal
	timer := time.NewTimer(time.Millisecond * 2000)

	//Requesting cab calls this elevator's unfinished cabcalls from other elevators
	sendMessageToElev(&msgID, message.ALL_ELEV, elevatorID, currentFloor, message.REQUEST_CAB_CALL, order.CreateOrder(currentFloor, elevio.BT_Cab), send_broadcast_message)
	for {
		select {
		//Button is polled from panel
		case newButton := <-drv_button:
			distributeOrderToElev(database, &msgID, currentFloor, elevatorID, order.CreateOrder(newButton.Floor, newButton.Button), newAliveMap, elevOnFloor, new_order, send_broadcast_message)

		//Receive broadcast message
		case newMessage := <-receive_broadcast_message:
			//Checking if the message is meant for this elevator
			if newMessage.FromID == elevatorID || (newMessage.ToID != elevatorID && newMessage.ToID != message.ALL_ELEV) {
				break
			}
			switch newMessage.MsgType {
			case message.ADD_ORDER:
				order.AddOrderToDatabase(database, elevatorID, newMessage.Orders[0])
				sendMessageToElev(&newMessage.MsgID, newMessage.FromID, elevatorID, currentFloor, message.ACK, newMessage.Orders[0], send_broadcast_message)
				new_order <- newMessage.Orders[0]
				setLightOnAllOtherElevators(database, &msgID, currentFloor, elevatorID, newAliveMap, newMessage.Orders[0], send_broadcast_message, true)

			case message.ACK:
				send_broadcast_message <- newMessage

			case message.I_AM_ALIVE:
				receive_alive <- newMessage.FromID
				if elevTaken[newMessage.FromID] {
					elevTaken[newMessage.FromID] = false
				}
				database = order.UpdateDatabaseQueue(database, newMessage.FromID, order.CreateQueueFromOrders(order.NumFloors, newMessage.Orders), order.NumFloors)

			case message.LIGHT_ON:
				set_lights <- order.SetLight(newMessage.Orders[0], true)
				sendMessageToElev(&newMessage.MsgID, newMessage.FromID, elevatorID, currentFloor, message.ACK, newMessage.Orders[0], send_broadcast_message)

			case message.LIGHT_OFF:
				set_lights <- order.SetLight(newMessage.Orders[0], false)
				sendMessageToElev(&newMessage.MsgID, newMessage.FromID, elevatorID, currentFloor, message.ACK, newMessage.Orders[0], send_broadcast_message)
				elevOnFloor[newMessage.FromID] = newMessage.Orders[0].Floor

			case message.REQUEST_CAB_CALL:
				if checkResponsibilityForCabcalls(newAliveMap, elevatorID) {
					cabCalls := order.CreateQueueFromOrders(order.NumFloors, order.GetCabcallsFromDatabase(database, newMessage.FromID))
					send_broadcast_message <- message.CreateMessageFromQueue(newMessage.MsgID, newMessage.FromID, elevatorID, currentFloor, message.REQUEST_CAB_CALL_ACK, cabCalls)
				}

			case message.REQUEST_CAB_CALL_ACK:
				//Sending received cab calls to fsm
				for _, currentOrder := range newMessage.Orders {
					order.AddOrderToDatabase(database, elevatorID, currentOrder)
					new_order <- currentOrder
				}
				//sending ack to messageServer to stop requesting for cab calls
				send_broadcast_message <- newMessage
			}

		//This elevator has finished an order
		case newOrderExecuted := <-order_finished:
			setLightOnAllOtherElevators(database, &msgID, currentFloor, elevatorID, newAliveMap, newOrderExecuted, send_broadcast_message, false)
			order.RemoveOrderFromDatabase(database, elevatorID, newOrderExecuted)
			currentFloor = newOrderExecuted.Floor
			elevOnFloor[elevatorID] = currentFloor

		//Received new updated map of elevators alive
		case newAliveMap = <-alive_map:
			fmt.Println(newAliveMap)
			//Checking if any elevators have died since the last time we received alive information
			for elevID, queue := range database {
				if !newAliveMap[elevID] && !elevTaken[elevID] {
					elevTaken[elevID] = true
					distributeFallenHallOrders(database, queue, &msgID, currentFloor, elevatorID, newAliveMap, elevOnFloor, new_order, send_broadcast_message)
				}
			}

		// IAmAlive will be sent every 250 ms
		case <-timer.C:
			timer.Reset(time.Millisecond * 250)
			supervisor.SendIAmAlive(msgID, elevatorID, currentFloor, send_broadcast_message, database)

		//Will trigger if single elevator in fsm is in moving state for more than 10 seconds
		case checkEngine := <-engine_failed:
			if checkEngine {
				timer.Stop()
			} else {
				timer.Reset(0)
			}
		}

		if msgID >= 5000 {
			msgID = 0
		}
	}
}
