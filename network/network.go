package network

import (
	"net"
	"log"
)

const (
	LOCAL_SERV_ADDR = "localhost"
	LOCAL_SERV_PORT = ":50123"
)

// Start the server with the given address and port
func startServer(){
	ln, err := net.Listen("tcp", LOCAL_SERV_PORT)
	if err != nil{
		log.Print("Failed to connect the server.")
	} else {
		log.Print("Server started on" + LOCAL_SERV_ADDR + LOCAL_SERV_PORT)
	}

	for {
		conn, _ := ln.Accept()
		go acceptBlock(conn)
	}
}

// Destring block data, check block ID against known block list
// Decide to accept or discard
func acceptBlock(conn net.Conn) {
	// Destringify the incoming data using function from block package
	// blockData := block.DestringifyBlockData(conn)
	/*
	 TODO: implement check incoming blocks
	and decide whether they should be dropped or accepted
	*/
}