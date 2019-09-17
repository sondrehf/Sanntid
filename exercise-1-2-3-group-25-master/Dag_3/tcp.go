package main 

import (
	"fmt"
	"net"
	"time"
)

func tcp_client(finished chan<- bool){

	port := "33546" 
	var network string = "tcp"
	broadcast2 := "10.100.23.242"
	tcpAdr, err1 := net.ResolveTCPAddr(network, broadcast2+":"+port)

	sock, err2 := net.DialTCP(network, nil, tcpAdr)
	defer sock.Close()

	fmt.Println("error 1:", err1,"\nerror 2:", err2)
	fmt.Println("\nSock 1:", sock)

	//sock.Write([]byte("Connect to: 10.100.23.254:20005"+"\x00"))
	sock.Write([]byte("ABCDEFGHIJKLMNOPQRST"+"\x00"))
	//time.Sleep(time.Second)
	reply := make([]byte, 1024)
 	n, _ := sock.Read(reply)
	fmt.Println(string(reply[:n]))
	n, _ = sock.Read(reply)
	fmt.Println(string(reply[:n]))

	finished <- true
}


func tcp_server(finished chan<- bool){
	
	port := "20005"
	var network string = "tcp"
	
	tcpAdr, err1 := net.ResolveTCPAddr(network, ":"+port) //our ip
	sock, err2 := net.ListenTCP(network, tcpAdr)

	fmt.Println("error 1s:", err1,"\nerror 2s:", err2)
	defer sock.Close()

	//for {
		reply := make([]byte, 1024)

		con, err3 := sock.AcceptTCP()
		fmt.Println("err3: ", err3, "n: ", con)
		p, _ := con.Read(reply)
		fmt.Println(string(reply[:p]))
		time.Sleep(time.Second)
		finished <- true


	//}

	//time.Sleep(time.Second)
}


func main(){
	finished := make(chan bool)

	//go tcp_server(finished)
	//time.Sleep(time.Second*5)

	go tcp_client(finished)

	//<- finished 
	<- finished
	

}