package fsm

import (
	"fmt"
	"../elevio"
	"../order"
)

// STATE for each elevator. Four possible states which is used in FSMFunc
type STATE int

const ( 
	INIT      STATE = iota 
	IDLE                   
	MOVING                 
	DOOR_OPEN              
)

// Elev struct is used to update the state machine in each elevator
type Elev struct {
	State        STATE
	Dir          elevio.MotorDirection
	Floor        int
	Queue        [][]bool
	CurrentOrder order.Order
}

//chooseDir returns the direction depending on the elevators queue
func chooseDir(elevator Elev) elevio.MotorDirection {
	switch elevator.Dir {
	case elevio.MD_Up:
		return elevator.Dir
	case elevio.MD_Down:
		return elevator.Dir
	case elevio.MD_Stop:
		if ordersBelow(elevator) {
			elevator.Dir = elevio.MD_Down
		}
		if ordersAbove(elevator) {
			elevator.Dir = elevio.MD_Up
		}
	}
	return elevator.Dir
}

//shouldElevatorStop checks if the elevator should stop at a given floor
func shouldElevatorStop(elevator Elev) bool {
	switch elevator.Dir {
	case elevio.MD_Up:
		return elevator.Queue[elevator.Floor][0] || elevator.Queue[elevator.Floor][2] || !ordersAbove(elevator)
	case elevio.MD_Down:
		return elevator.Queue[elevator.Floor][1] || elevator.Queue[elevator.Floor][2] || !ordersBelow(elevator)

	case elevio.MD_Stop:
	}
	return false
}

//ordersBelow checks if there is any active orders below the current floor
func ordersBelow(elevator Elev) bool {
	for i := elevator.Floor - 1; i >= 0; i-- {
		for j := 0; j < order.NumBtns; j++ {
			if elevator.Queue[i][j] {
				return true
			}
		}
	}
	return false
}

//ordersAbove checks if there is any active orders above the current floor
func ordersAbove(elevator Elev) bool {
	for i := elevator.Floor + 1; i < order.NumFloors; i++ {
		for j := 0; j < order.NumBtns; j++ {
			if elevator.Queue[i][j] {
				return true
			}
		}
	}
	return false
}

//deleteFloor removes all active orders at the current floor in the local queue  when an order is executed
func deleteFloor(elevator Elev) {
	for i, btn := range elevator.Queue[elevator.Floor] {
		if btn {
			order.RemoveFromQueue(elevator.Queue, order.CreateOrder(elevator.Floor, elevio.ButtonType(i)))
			elevio.SetButtonLamp(elevio.ButtonType(i), elevator.Floor, false)
		}
	}
}

//drvOrderFinished triggers the order_finished channel when an floor is reached
func drvOrderFinished(elevator Elev, order_finished chan order.Order){
	for i, btn := range elevator.Queue[elevator.Floor] {
		if btn {
			order_finished <- order.CreateOrder(elevator.Floor, elevio.ButtonType(i))
		}
	}
}

//printQueue prints the updated queue when the state machine reciews a new order
func printQueue(elevator Elev) {
	fmt.Println("New order registered :")
	fmt.Println("Floor \t | Up\t| Down\t| Cab")
	for i := order.NumFloors - 1; i >= 0; i-- {
		fmt.Printf("Floor: %d |", i)
		fmt.Println(elevator.Queue[i])
	}
	fmt.Println()
}
