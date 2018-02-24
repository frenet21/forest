package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"

	"./block"
	"./pool"
)

func main() {
	fmt.Println("Initializing....")

	pool.GenesisPool()

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	fmt.Println("Done.\n")

	for true {
		fmt.Println("Please enter a message:")
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		fmt.Println("")
		if err != nil {
			panic(err)
		}

		fmt.Println("Creating message block...")
		transmission := block.CreateBlock(message, &publicKey)
		fmt.Println("Message block created.\n")
		fmt.Println("")

		fmt.Println("Stringifying block....")
		str := block.StringifyBlock(transmission)
		fmt.Println("Done.")
		fmt.Println("")
		fmt.Printf("Stringified block (base64): %s\n\n", base64.URLEncoding.EncodeToString([]byte(str)))
		fmt.Println("")

		fmt.Println("Destringifying block...")
		reception := block.DestringifyBlock(str)
		fmt.Println("Block recreated.\n")
		fmt.Println("")

		fmt.Println("Attempting block decryption...")
		output, err := block.AttemptDecrypt(reception, privateKey)
		if err != nil {
			panic(err)
		}
		fmt.Println("Block decrypted.\n")
		fmt.Println("")

		fmt.Printf("Original message was: %s\n", message)
		fmt.Printf("Decrypted message was: %s\n", output)
		if message == output {
			fmt.Println("Messages match.")
		} else {
			fmt.Println("Messages differ.")
		}
		fmt.Println("")
	}
}
