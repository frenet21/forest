package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"../block"
	"../pool"
)

const (
	ID_LIST_PATH       = "BLOCK_ID_LIST.txt"
	KNOWN_CLIENTS_PATH = "KNOWN_CLIENTS.txt"
	LOCAL_SERV_ADDR    = "localhost"
	LOCAL_SERV_PORT    = ":50123"
	RECEIVER_PORT      = ":50123"
)

/*
startServer() 	- 	Starts the listener on the address. Starts acceptBlock() as a goroutine.
acceptBlock() 	- 	Receives blocks, destrings, accepts/drops based on the known hash list.
fowardBlock() 	- 	Loops the client list, firing off sendBlock(), then sends block to the frontend.
					Clients sending messages will simply pass them into this function.
sendBlock()		-	Sends a given block to an address (IP:PORT)
*/

// Start the server with the given address and port
func startServer(done chan bool) {
	ln, err := net.Listen("tcp", LOCAL_SERV_PORT)
	if err != nil {
		log.Print("[NET - SERVER] Failed to start the server.")
	} else {
		log.Print("[NET - SERVER] Server started on" + LOCAL_SERV_ADDR + LOCAL_SERV_PORT)
	}

	// Server start is completed
	done <- true

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
	var blockbytes []byte
	var len int
	if n, err := conn.Read(blockbytes); err != nil {
		panic(err)
	} else {
		len = n
	}
	Block := block.DestringifyBlock(string(blockbytes[:len]))
	// Select the block ID from the Block
	blockID := string(Block.ID[:64])

	// Check if the known hash text file exists
	_, err := os.Stat(ID_LIST_PATH)
	if err != nil {
		file, err := os.Create(ID_LIST_PATH)
		if err != nil {
			log.Print("[NET - ACCEPTOR] Failed to create file.")
		} else {
			log.Print("[NET - ACCEPTOR] A new ID list file was created.")
		}
	}

	// Search the known_hash.txt file for the blockID
	f, err := os.Open(ID_LIST_PATH)
	if err != nil {
		log.Print("[NET - ACCEPTOR] Could not open the file.")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// If it is there, discard the Block and stop this function
		if strings.Contains(scanner.Text(), blockID) {
			log.Print("[NET - ACCEPTOR] Received block is known. Discarding...")
			f.Close()
		} else {
			log.Print("[NET - ACCEPTOR] New block ID identified. Adding to list...")
			f, err := os.OpenFile(ID_LIST_PATH, os.O_APPEND|os.O_WRONLY, 0644)
			n, err := f.WriteString(blockID)
			f.Close()

			// Pass block to forwardBlock function
			forwardBlock(Block)
		}
	}
}

// acceptBlock function passes OK'd blocks here
// These blocks are sent to everyone on the client list
// And the decryption attempt function is called here
func forwardBlock(blk block.Block) {
	// Give this block to the blockpool
	pool.ReceiveBlockHash(string(blk.ID[:64]))

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
		// Send block and the read address to sendBlock function
		sendBlock(blk, sendAddress)
	}
	if err != nil {
		log.Print("[NET - FORWARDER] Failed to read opened client list.")
	}

	/*
		TODO: Loop through files or entries of private keys
	*/

	// Attempt a decryption of the received block after passing to known clients
	message, err := block.AttemptDecrypt(blk, priKey)
	if err != nil {
		log.Print("[NET - FORWARDER] Decryption failed. Discarding block.")
	} else {
		log.Print("[NET - FORWARDER] Decryption success. Sending message to frontend.")
		log.Print("[NET - FORWARDER] \n[BEGIN DECRYPTED MESSAGE]\n" + message + "\n[END DECRYPTED MESSAGE]")
		/*
			TODO: Send block to frontend function
		*/
	}
}

// Send a block to a given address.
// Important note: sendAddress should be stored as IP:PORT
func sendBlock(blk block.Block, sendAddress string) {
	log.Print("[NET - SENDER] Dialing " + sendAddress + "... ")
	// Translate the address string to a dialable TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", sendAddress)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	// Skip attempted send if address dial fails.
	if err != nil {
		log.Print("[NET - SENDER] Failed dialing " + sendAddress + ". Skipping.")
	}

	// Send block to socket
	conn.Write([]byte(block.StringifyBlock(blk)))
}
