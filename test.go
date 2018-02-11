package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
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
		if err != nil {
			panic(err)
		}

		transmission := block.CreateBlock(message, &publicKey)

		str := block.StringifyBlock(transmission)

		reception := block.DestringifyBlock(str)

		output, err := block.AttemptDecrypt(reception, privateKey)
		if err != nil {
			panic(err)
		}

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
