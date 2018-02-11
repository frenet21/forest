package network

import (
	"net"
	"log"
)

const (
	LOCAL_SERV_ADDR = "localhost"
	LOCAL_SERV_PORT = ":50123"
)

func startServer(){
	ln, err := net.Listen("tcp", LOCAL_SERV_PORT)
	if err != nil{
		log.Print("Failed to connect the server.")
	} else {
		log.Print("Server started on" + LOCAL_SERV_ADDR + LOCAL_SERV_PORT)
	}

	for{
		conn, _ := ln.Accept()
		go listen(conn)
	}
}

func Listen(){
	
}