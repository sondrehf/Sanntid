package supervisor

import (
	"time"
	"../message"
)

//SendIAmAlive sends an IAmAlive message including the local elevators queue to the other elevators
func SendIAmAlive(msgID int, fromID int, currentFloor int, send_broadcast_message chan<- message.Message, database [][][]bool) {
	send_broadcast_message <- message.CreateMessageFromQueue(msgID, message.ALL_ELEV, fromID, currentFloor, message.I_AM_ALIVE, database[fromID])
}

//SetElevatorsAliveINIT called as a goroutine that sets the first active elevator alive.  
func SetElevatorsAliveInit(receive_alive <-chan int, alive_map chan<- map[int]bool) {
	alive := make(map[int]bool)
	timer := time.NewTimer(time.Second * 2)
	for {
		select {
		case index := <-receive_alive:
			if len(alive) <= index {
				alive = appendElevators(alive, index-len(alive)+1)
			}
			alive[index] = true
		case <-timer.C:
			alive_map <- alive
			return
		}
	}
}

//SetElevatorsAlive called as a goroutine that keep track of which elevators that are alive 
func SetElevatorsAlive(elevatorID int, receive_alive <-chan int, alive_map chan<- map[int]bool) {
	alive := createAlive(elevatorID)
	alive[elevatorID] = true
	counter := make(map[int]int)
	timer := time.NewTimer(900 * time.Millisecond)
	for {
		select {
		case index := <-receive_alive:
			if alive[index] {
				counter[index]++
			}
			if len(alive) <= index && counter[index] > 0 {
				alive = appendElevators(alive, index-len(alive)+1)
			}
			alive[index] = true
		case <-timer.C:
			timer.Reset(1000 * time.Millisecond)
			alive_map <- alive
			alive[elevatorID] = true
			alive = createAlive(elevatorID)
			counter = make(map[int]int)
		}
	}
}



