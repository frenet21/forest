package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"./block"
	"./pool"
)

func main() {
	pool.GenesisPool()

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	message1 := "THIS IS a TeSt!1"
	message2 := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec ac lobortis leo, in malesuada ipsum. Sed et tincidunt urna, a consectetur risus. Curabitur non orci at justo egestas rhoncus. Nullam et ultricies arcu. Ut pharetra dolor ac viverra sollicitudin. Quisque bibendum mi a sem aliquam, eget bibendum orci finibus. Proin ut neque vitae ligula interdum laoreet. Ut vel metus faucibus, rutrum lacus nec, gravida dui. Integer lobortis malesuada sodales. Pellentesque scelerisque est pulvinar ante cursus aliquet."

	transmission1 := block.CreateBlock(message1, &publicKey)
	transmission2 := block.CreateBlock(message2, &publicKey)

	output1, err := block.AttemptDecrypt(transmission1, privateKey)
	if err != nil {
		panic(err)
	}
	output2, err := block.AttemptDecrypt(transmission2, privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Preencrypt: %s. Postencrypt: %s. Matches: %t", message1, output1, message1 == output1)
	fmt.Printf("Preencrypt: %s. Postencrypt: %s. Matches: %t", message2, output2, message2 == output2)
}
