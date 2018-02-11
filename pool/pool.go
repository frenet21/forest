// pool
package main 

import (
	"golang.org/x/crypto/sha3"
	"fmt"
)

type encryptedMessageAndHash struct {
	encryptedMessage string
	blockpool [1000]byte
}

// Selects a block parent based on the encrypted message
func selectParentHash(info encryptedMessageAndHash) [64]byte {
	// TODO: Connect this to the blockpool
	var out [64]byte
	copy(out[:], sha3.New512().Sum([]byte(info.encryptedMessage)))
	return out
}

func main() {
	var a [1000]byte
	info := encryptedMessageAndHash{encryptedMessage: "Bobby", blockpool: a}
	out := selectParentHash(info)
	fmt.Printf("%x", out)
}
