package network

import (
	"net"
	"log"
	"os"
	"encoding/gob"
	"forest/block"
)

const (
	ID_LIST_PATH = "/network/ID_LIST.txt"
	LOCAL_SERV_ADDR = "localhost"
	LOCAL_SERV_PORT = ":50123"
)

type IDStore struct {
	ID string
}

// Start the server with the given address and port
func startServer(){
	ln, err := net.Listen("tcp", LOCAL_SERV_PORT)
	if err != nil{
		log.Print("Failed to start the server.")
	} else {
		log.Print("Server started on" + LOCAL_SERV_ADDR + LOCAL_SERV_PORT)
	}

	for {
		conn, _ := ln.Accept()
		// Start function goroutine to accept connection
		go acceptBlock(conn)
	}
}

// Destring block, check block ID against known block list
// Decide to accept or discard
// Accepted block is returned to the frontend for storing and later viewing

func acceptBlock(conn net.Conn) Block {
	// Destringify the "conn" string with function from package 'block'
	Block := block.DestringifyBlock(conn)
	// Select the block ID from the Block
	blockID := Block.ID

	// Check if the known_hash.txt file exists
	if _, err := os.Stat(ID_LIST_PATH) {
		// If not, create it...
		file, err := os.Create(ID_LIST_PATH)
		if err != nil {
			log.Print("Failed to create file.")
			return 0
		} else {
			log.Print("A new ID list file was created.")
		}
	}

	// Search the known_hash.txt file for the blockID
	f, err := os.Open(path)
	if err != nil {
		log.Print("Could not open the file.")
		return 0
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// If it is there, discard the Block and stop this function
		if strings.Contains(scanner.Text(), blockID) {
			log.Print("Received block is known. Discarding...")
			/*
			TODO: Discard block function
			*/
			return line, nil
		}
		line++
	}
	
	// If it is not, return the Block to the frontend
	/*
	TODO: Pass on block for storage on frontend
	*/
}


// After acceptBlock accepts block, and stored for later viewing
func receiveBlock() {

}

func sendBlock() {

}
