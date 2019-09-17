package fsm

import (
	"fmt"
	"time"
	"../elevio"
	"../order"
)

//FSMFunc called as a goroutine, runs single elevator and updates orderHandler.
func FsmFunc(new_order <-chan order.Order, order_finished chan order.Order, set_lights <-chan order.Light, lights_finished chan<- order.Light, engine_failed chan<- bool) {
	elevator := Elev{
		State: INIT,
		Dir:   elevio.MD_Stop,
		Floor: elevio.GetFloor(),
		Queue: order.CreateEmptyQueue(order.NumFloors),
	}

	var dir elevio.MotorDirection = elevio.MD_Down

	drv_floors := make(chan int)
	engineTimedOut := time.NewTimer(10 * time.Second)
	doorTimedOut := time.NewTimer(0)

	engineTimedOut.Stop()

	go elevio.PollFloorSensor(drv_floors)

	if elevator.State == INIT {
		dir = elevio.MD_Down
		elevio.SetMotorDirection(dir)
		initFloor := <-drv_floors
		elevator.Floor = initFloor
		elevio.SetFloorIndicator(initFloor)
		fmt.Printf("Current floor:%+v\n", initFloor)
		dir = elevio.MD_Stop
		elevio.SetMotorDirection(dir)
		elevator.State = IDLE

		order_finished <- order.CreateOrder(initFloor, elevio.BT_Cab)
	}
	//The elevator has 5 possible events: new order, floor reached, door closing, engine failed and set lights. 
	for {
		select {
		case myOrder := <-new_order: //NEW ORDER RECEIVED
			order.AddToQueue(elevator.Queue, myOrder)
			elevio.SetButtonLamp(myOrder.Button, myOrder.Floor, true)
			printQueue(elevator)
			switch elevator.State {
			case IDLE:
				elevator.Dir = chooseDir(elevator)
				if elevator.Dir == elevio.MD_Stop {
					elevio.SetDoorOpenLamp(true)
					doorTimedOut.Reset(3 * time.Second)
					deleteFloor(elevator)
					go func() { order_finished <- myOrder }()
					elevator.State = DOOR_OPEN

				} else {
					elevio.SetMotorDirection(elevator.Dir)
					engineTimedOut.Reset(10 * time.Second)
					elevator.State = MOVING
				}
			case MOVING:
			case DOOR_OPEN:
				if myOrder.Floor == elevator.Floor {
					doorTimedOut.Reset(3 * time.Second)
					deleteFloor(elevator)
					go func() { order_finished <- myOrder }()
					elevio.SetMotorDirection(elevator.Dir)
				}
			}
		case a := <-drv_floors: //FLOOR REACHED
			elevator.Floor = a
			elevio.SetFloorIndicator(elevator.Floor)
			if shouldElevatorStop(elevator) {
				elevator.Dir = elevio.MD_Stop
				elevio.SetDoorOpenLamp(true)
				engineTimedOut.Stop()
				elevator.State = DOOR_OPEN
				doorTimedOut.Reset(3 * time.Second)
				elevio.SetMotorDirection(elevator.Dir)
				drvOrderFinished(elevator, order_finished)
				deleteFloor(elevator)
			} else if elevator.State == MOVING {
				engineTimedOut.Reset(10 * time.Second)
			}
		case <-doorTimedOut.C: //DOOR TIMED OUT
			engine_failed <- false
			elevio.SetDoorOpenLamp(false)
			elevator.Dir = chooseDir(elevator)
			if elevator.Dir == elevio.MD_Stop {
				elevator.State = IDLE
				engineTimedOut.Stop()
			} else {
				elevator.State = MOVING
				engineTimedOut.Reset(10 * time.Second)
				elevio.SetMotorDirection(elevator.Dir)
			}

		case <-engineTimedOut.C:
			fmt.Println("Engine error! Do something clever ")
			engine_failed <- true

		case b := <-set_lights:
			elevio.SetButtonLamp(b.Orders.Button, b.Orders.Floor, b.SetLight)
		}
	}
}
