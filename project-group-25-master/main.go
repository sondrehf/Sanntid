package main

import (
	"flag"
	"./elevio"
	"./message"
	"./order"
	"./orderHandler"
)

func main() {
	var port string
	var elevatorID int
	flag.IntVar(&elevatorID, "id", 0, "Choose id (0 is standard)")
	flag.StringVar(&port, "port", "15657", "Choose port")
	flag.Parse()
	elevio.Init("localhost:", port, order.NumFloors)

	//Creating channels
	send_broadcast_message := make(chan message.Message)
	broadcast_message := make(chan message.Message)
	receive_broadcast_message := make(chan message.Message)
	order_finished := make(chan order.Order)
	new_order := make(chan order.Order)
	drv_buttons := make(chan elevio.ButtonEvent)
	receive_alive := make(chan int)
	alive_list := make(chan map[int]bool)
	engine_fail := make(chan bool)
	set_lights := make(chan order.Light)
	lights_finished := make(chan order.Light)
	
	orderHandler.OrderServer(elevatorID, drv_buttons, receive_broadcast_message, order_finished, alive_list, send_broadcast_message, broadcast_message, new_order, receive_alive, set_lights, lights_finished, engine_fail)
	select {}
}
