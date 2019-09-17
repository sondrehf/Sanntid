package message

import (
	"time"
	"../elevio"
	"../order"
)

type MessageType int

const (
	ADD_ORDER = iota
	ACK
	I_AM_ALIVE
	LIGHT_ON
	LIGHT_OFF
	REQUEST_CAB_CALL
	REQUEST_CAB_CALL_ACK
)

type Message struct {
	MsgID        int
	ToID         int
	FromID       int
	CurrentFloor int
	MsgType      MessageType
	Orders       []order.Order
}

type ElevID int
const (
	ALL_ELEV = -1
)

//CreateMessageFromOrders returns a message struct in the right format from one or multiple orders
func CreateMessageFromOrders(msgID int, toID int, fromID int, currentFloor int, msgType MessageType, orders ...order.Order) Message {
	return Message{msgID, toID, fromID, currentFloor, msgType, orders}
}

//CreateMessageFromOrders returns a message struct in the right format from an elevator's queue. Used when retrieving cab calls
func CreateMessageFromQueue(msgID int, toID int, fromID int, currentFloor int, msgType MessageType, queue [][]bool) Message {
	orders := []order.Order{}
	for floor, myButtons := range queue {
		for myButtonType, active := range myButtons {
			if active {
				orders = append(orders, order.CreateOrder(floor, elevio.ButtonType(myButtonType)))
			}
		}
	}
	return Message{msgID, toID, fromID, currentFloor, msgType, orders}
}

//MessageServer called as a goroutine that transmits all active messages repetitious (maximum 10 times), 
//and stop sending them when receiving an acknowledgement from the receiver. 
func MessageServer(elevatorID int, send_broadcast_message <-chan Message, broadcast_message chan<- Message) {
	activeMessages := make(map[int]Message) //contains messages that have not been acked
	counters := make(map[int]int)           // counter for each message

	//timer used for repeatedly sending unAcked messages
	timer := time.NewTimer(time.Millisecond * 50)

	for {
		select {
		case newMessage := <-send_broadcast_message: //receiving a message from orderServer
			//Ack messages received from other elevators will stop this elevator from sending new messages
			if (newMessage.MsgType == ACK || newMessage.MsgType == REQUEST_CAB_CALL_ACK) && newMessage.FromID != elevatorID {
				delete(activeMessages, newMessage.MsgID)
				delete(counters, newMessage.MsgID)

			//Ack and IAmAlive messages created by this elevator itself will be broadcasted
			} else if (newMessage.MsgType == ACK && newMessage.FromID == elevatorID) || (newMessage.MsgType == REQUEST_CAB_CALL_ACK && newMessage.FromID == elevatorID) || newMessage.MsgType == I_AM_ALIVE {
				broadcast_message <- newMessage

			//New messages from this elevator
			} else {
				activeMessages[newMessage.MsgID] = newMessage
				counters[newMessage.MsgID] = 1
				broadcast_message <- newMessage
			}

		//resending unAcked messages
		case <-timer.C: 
			timer.Reset(time.Millisecond * 50)
			for _, currentMessage := range activeMessages {
				if counters[currentMessage.MsgID] > 10 {
					delete(activeMessages, currentMessage.MsgID)
					delete(counters, currentMessage.MsgID)
				} else {
					counters[currentMessage.MsgID]++
					broadcast_message <- currentMessage
				}
			}
		}
	}
}


