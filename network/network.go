package network

import (
	"net"
	"log"
	"os"
	"encoding/gob"
	"forest/block"
)

const (
	ID_LIST_PATH = "/network/BLOCK_ID_LIST.txt"
	KNOWN_CLIENTS_PATH = "/network/KNOWN_CLIENTS.txt"
	LOCAL_SERV_ADDR = "localhost"
	LOCAL_SERV_PORT = ":50123"
)

// Start the server with the given address and port
func startServer(done chan bool){
	ln, err := net.Listen("tcp", LOCAL_SERV_PORT)
	if err != nil{
		log.Print("[NET - SERVER] Failed to start the server.")
	} else {
		log.Print("[NET - SERVER] Server started on" + LOCAL_SERV_ADDR + LOCAL_SERV_PORT)
	}

	// Server start is completed
	done<-true

	// Start function goroutine to accept connection
	for {
		conn, _ := ln.Accept()
		go acceptBlock(conn)
	}
}

// Destring block, check block ID against known block list
// Decide to accept or discard
// Accepted block is passed to forwardBlock function
func acceptBlock(conn net.Conn) {
	// Destringify the "conn" string with function from package 'block'
	Block := block.DestringifyBlock(conn)
	// Select the block ID from the Block
	blockID := Block.ID

	// Check if the known hash text file exists
	if _, err := os.Stat(ID_LIST_PATH) {
		// If not, create it...
		file, err := os.Create(ID_LIST_PATH)
		if err != nil {
			log.Print("[NET - ACCEPTOR] Failed to create file.")
			return 0
		} else {
			log.Print("[NET - ACCEPTOR] A new ID list file was created.")
		}
	}

	// Search the known_hash.txt file for the blockID
	f, err := os.Open(path)
	if err != nil {
		log.Print("[NET - ACCEPTOR] Could not open the file.")
		return 0
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// If it is there, discard the Block and stop this function
		if strings.Contains(scanner.Text(), blockID) {
			log.Print("[NET - ACCEPTOR] Received block is known. Discarding...")
			f.close()
			return 0
		}
		// If it is not, add the block ID to the list 
		else {
			f, err := os.OpenFile(ID_LIST_PATH, os.O_APPEND|os.O_WRONLY, 0644) 
			n, err := f.WriteString(blockID) 
			f.Close()

			// Pass block to forwardBlock function
			forwardBlock(Block)
			return 0
		}
		line++
	}
}


// acceptBlock function passes OK'd blocks here
// These blocks are sent to everyone on the client list
// And the decryption attempt function is called here
func forwardBlock(block Block) {
	// Open the known client list path
	file, err := os.Open(KNOWN_CLIENTS_PATH)
	if err != nil {
		log.Print("[NET - FORWARDER] Failed to open client list.")
	}
	defer file.Close()

	// Scan each line of the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Send off to each client on the list
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		sendAddress := scanner.Text()
		sendBlock(block, sendAddress)
	}
	if err != nil {
		log.Print("[NET - FORWARDER] Failed to read opened client list.")
	}

	// Attempt a decryption of the received block after passing to known clients
	message, err := block.AttemptDecrypt(block, priKey)
	if err != nil {
		log.Print("[NET - FORWARDER] Decryption failed. Discarding block.")
		return 0
	} else {
		log.Print("[NET - FORWARDER] Decryption success. Sending message to frontend.")
		/*
		TODO: Send block to frontend function
		*/
	}
}