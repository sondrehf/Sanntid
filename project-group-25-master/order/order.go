package order

import (
	"../elevio"
)

//const variable
const (
	NumFloors = 4
	NumBtns   = 3
)

//Order renaming elevio.ButtonEvent to Order
type Order elevio.ButtonEvent

//Light is a struct used to turn lights on and off for one order
type Light struct {
	SetLight bool
	Orders   Order
}

//CreateOrder takes the arguments floor and a buttontype, and returns an order struct
func CreateOrder(floor int, button elevio.ButtonType) Order {
	return Order{floor, button}
}

//SetLight returns an Light struct for the given order and situation.
func SetLight(orders Order, setLight bool) Light {
	return Light{setLight, orders}
}

//CreateEmptyQueue creates an empty queue at each elevator
func CreateEmptyQueue(numberOfFloors int) [][]bool {
	queue := [][]bool{}
	for i := 0; i < numberOfFloors; i++ {
		row := []bool{false, false, false}
		queue = append(queue, row)
	}
	return queue
}

//CreateQueueFromOrders creates an empty queue and adds all the active orders
func CreateQueueFromOrders(numberOfFloors int, orders []Order) [][]bool {
	queue := CreateEmptyQueue(numberOfFloors)
	for _, currentOrder := range orders {
		AddToQueue(queue, currentOrder)
	}
	return queue
}

//AddToQueue adds an order to the local queue
func AddToQueue(queue [][]bool, order Order) {
	queue[order.Floor][order.Button] = true
}

//RemoveFromQueue removes an order from the local queue
func RemoveFromQueue(queue [][]bool, order Order) {
	queue[order.Floor][order.Button] = false
}

//CreateEmptyDatabase creates a third dimensional array with orders for all active elevators
func CreateEmptyDatabase(elevatorID int, numberOfFloors int) [][][]bool {
	database := [][][]bool{}
	database = addNewElevatorsToDatabase(database, elevatorID+1, numberOfFloors)
	return database
}

//UpdateDatabaseQueue updates the local queue at each elevator
func UpdateDatabaseQueue(database [][][]bool, elevatorID int, queue [][]bool, numberOfFloors int) [][][]bool {
	numberOfElevators := len(database)
	if numberOfElevators <= elevatorID {
		numberOfNewElevators := elevatorID - (numberOfElevators - 1)
		database = addNewElevatorsToDatabase(database, numberOfNewElevators, numberOfFloors)
	}
	database[elevatorID] = queue
	return database
}

//AddOrderToDatabase adds an order to the database
func AddOrderToDatabase(database [][][]bool, elevatorID int, order Order) {
	database[elevatorID][order.Floor][order.Button] = true
}

//RemoveOrderFromDatabase removes an order from the database
func RemoveOrderFromDatabase(database [][][]bool, elevatorID int, order Order) {
	database[elevatorID][order.Floor][order.Button] = false
}

//GetCabcallsFromDatabase retrieves the cabcalls for a given elevator from the database
func GetCabcallsFromDatabase(database [][][]bool, elevID int) []Order {
	cabCalls := []Order{}
	if elevID < len(database) {
		for floor, myButtons := range database[elevID] {
			if myButtons[int(elevio.BT_Cab)] {
				cabCalls = append(cabCalls, CreateOrder(floor, elevio.BT_Cab))
			}
		}
	}
	return cabCalls
}

