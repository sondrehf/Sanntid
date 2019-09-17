package orderHandler

import (
	"fmt"
	"math"

	"../elevio"
	"../message"
	"../order"
	"../supervisor"
)

//initDatabaseAndAliveMap initialize the database and aliveMap when the program starts
func initDatabaseAndAliveMap(elevatorID int, alive_list chan map[int]bool, receive_alive chan int, receive_broadcast_message <-chan message.Message) ([][][]bool, map[int]bool) {
	database := [][][]bool{}
	newAliveMap := make(map[int]bool)
	go supervisor.SetElevatorsAliveInit(receive_alive, alive_list)

	initFinished := false
	for {
		select {
		case newAliveMap = <-alive_list:
			newAliveMap[elevatorID] = true
			fmt.Println(newAliveMap)
			database = order.CreateEmptyDatabase(elevatorID, order.NumFloors)
			initFinished = true
		case newMessage := <-receive_broadcast_message:
			if newMessage.MsgType == message.I_AM_ALIVE {
				receive_alive <- newMessage.FromID
			}
		}
		if initFinished {
			break
		}
	}
	return database, newAliveMap
}

//distributeOrderToElev determines which off the elevators that gets the order based on the button and our cost function
func distributeOrderToElev(
	database [][][]bool,
	msgID *int, currentFloor int,
	elevID int, newOrder order.Order,
	aliveList map[int]bool,
	elevOnFloor map[int]int,
	new_order chan<- order.Order,
	send_broadcast_message chan<- message.Message) {

	if newOrder.Button == elevio.BT_Cab {
		order.AddOrderToDatabase(database, elevID, newOrder)
		new_order <- newOrder
	} else if !isOrderUnderHandling(database, aliveList, order.CreateOrder(newOrder.Floor, newOrder.Button)) {
		toID := costFunction(database, aliveList, elevOnFloor, newOrder)
		if toID == elevID {
			order.AddOrderToDatabase(database, elevID, newOrder)
			setLightOnAllOtherElevators(database, msgID, currentFloor, elevID, aliveList, newOrder, send_broadcast_message, true)
			new_order <- newOrder
		} else {
			sendMessageToElev(msgID, toID, elevID, currentFloor, message.ADD_ORDER, newOrder, send_broadcast_message)
		}
	}

}

//isOrderUnderHandling checks if the new order already is under distribution
func isOrderUnderHandling(database [][][]bool, aliveList map[int]bool, newOrder order.Order) bool {
	status := false
	for elevatorID, _ := range database {
		if aliveList[elevatorID] == false {
			continue
		} else {
			status = database[elevatorID][newOrder.Floor][newOrder.Button]
			if status {
				return status
			}
		}
	}
	return status
}

//costFunction returns the elevator that is best suited to take the new hall order.
func costFunction(database [][][]bool, aliveList map[int]bool, elevOnFloor map[int]int, newOrder order.Order) int {
	id := -1
	numberOfOrders := order.NumFloors * order.NumBtns
	floorDistance := order.NumFloors
	for elevatorID, queue := range database {
		if aliveList[elevatorID] == false {
			continue
		}
		currentNumberOfOrders := 0
		for _, myButtons := range queue {
			for _, active := range myButtons {
				if active {
					currentNumberOfOrders++
				}
			}
		}

		currentFloorDistance := int(math.Abs(float64(newOrder.Floor) - float64(elevOnFloor[elevatorID])))
		fmt.Println("Elevator id: ", elevatorID, " current floor dist: ", currentFloorDistance)
		if currentNumberOfOrders < numberOfOrders {
			numberOfOrders = currentNumberOfOrders
			id = elevatorID
		}
		if currentNumberOfOrders == 0 && currentFloorDistance < floorDistance {
			floorDistance = currentFloorDistance
			id = elevatorID
		}
	}
	return id
}

//sendMessageToElev sends an order message to the message server.
func sendMessageToElev(msgID *int, toID int, fromID int, currentFloor int, msgType message.MessageType, sendOrder order.Order, send_broadcast_message chan<- message.Message) {
	if msgType != message.ACK && msgType != message.REQUEST_CAB_CALL_ACK {
		(*msgID)++
	}
	send_broadcast_message <- message.CreateMessageFromOrders(*msgID, toID, fromID, currentFloor, msgType, sendOrder)
}

//setLightOnAllElevator turns the lights on or off at the other active elevators.
func setLightOnAllOtherElevators(database [][][]bool, msgID *int, currentFloor int, elevID int, aliveList map[int]bool, sendOrder order.Order, send_broadcast_message chan<- message.Message, light bool) {
	for elevatorID, _ := range database {
		if aliveList[elevatorID] == true && elevID != elevatorID { //if key does  not exist, it will return false
			if light == true {
				sendMessageToElev(msgID, elevatorID, elevID, currentFloor, message.LIGHT_ON, sendOrder, send_broadcast_message)
			} else {
				sendMessageToElev(msgID, elevatorID, elevID, currentFloor, message.LIGHT_OFF, sendOrder, send_broadcast_message)
			}
		}
	}
}

//distributeFallenHallOrders distributes any remaining hall orders from an elevator that goes offline
func distributeFallenHallOrders(database [][][]bool, queue [][]bool, msgID *int, currentFloor int, elevID int, aliveList map[int]bool, elevOnFloor map[int]int, new_order chan<- order.Order, send_broadcast_message chan<- message.Message) {
	for floor, button := range queue {
		for btnType, active := range button {
			if active && btnType != int(elevio.BT_Cab) {
				distributeOrderToElev(database, msgID, currentFloor, elevID, order.CreateOrder(floor, elevio.ButtonType(btnType)), aliveList, elevOnFloor, new_order, send_broadcast_message)
			}
		}
	}
}

//chechResponsibilityForCabcalls determines which of the online elevators that have the responsibility to send the saved cabcalls when a
//new elevator returns online after a crash.
func checkResponsibilityForCabcalls(aliveMap map[int]bool, elevatorID int) bool {
	lowestID := len(aliveMap)
	for elevID, active := range aliveMap {
		if elevID < lowestID && active {
			lowestID = elevID
		}
	}
	return lowestID == elevatorID
}
