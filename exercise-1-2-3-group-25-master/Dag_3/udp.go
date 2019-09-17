package main 

import (
	"fmt"
	"net"
	"time"
)

func udp_sender(finished chan<- bool){

	var port string = "20005" 
	//var broadcast2 string = "10.100.23.242" //ip for master pc
	var network string = "udp"
	udpAdr, err1 := net.ResolveUDPAddr(network,/*broadcast2*/":"+port)
	sock, err2 := net.DialUDP(network, nil, udpAdr)

	fmt.Println("error 1:", err1,"\nerror 2:", err2)
	fmt.Println("\nSock 1:", sock)

	defer sock.Close()
	counter := 0
	for{
		sock.Write([]byte("Hello from group 25"))
		fmt.Println("\nnr msg sent: ", counter)
		time.Sleep(time.Second)
		counter++
		if counter >= 10{
			finished <- true
			return 
		}
	}
}

func udp_receiver(finished chan<- bool){
	buffer := make([]byte, 1024) //create buffer for reveived msg
	//var port2 string = "30000" //port for master pc 
	port := "20005" //student pc 
	var network string = "udp"

	udpAdr, err1 := net.ResolveUDPAddr(network, ":"+port)
	sock, err2 := net.ListenUDP(network, udpAdr)

	fmt.Println("error 1:", err1,"\nerror 2:", err2)

	defer sock.Close()
	counter := 0
	for {
		n, _ := sock.Read(buffer)
		fmt.Println(string(buffer[:n]))
		fmt.Println("\nnr msg received: ", counter)
		counter++
		if counter >= 10{
			finished <- true
			return 
		}
	}
}


/*func main(){
	finished := make(chan bool)

	go udp_sender(finished)
	go udp_receiver(finished)

	<- finished 
	<- finished
}*/