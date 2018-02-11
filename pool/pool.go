package pool

import (
	"golang.org/x/crypto/sha3"
	"fmt"
	"sort"
	"strings"
)

type encryptedMessageAndHash struct {
	encryptedMessage string
	blockpool []string
}

// Selects a block parent based on the encrypted message
func selectParentHash(info encryptedMessageAndHash) string {
	//Add hash of encrypted message to the end of the blockpool array
	hash := string(sha3.New512().Sum(byte[](info.encryptedMessage))[:64])
	blockpoolStrings := append(info.blockpool, hash)

	//Sorted blockpool strings
	fmt.Println(blockpoolStrings)
	sort.Strings(blockpoolStrings)
	fmt.Println(blockpoolStrings)

	var element int

	//Determine the element after the encrypted message
	for i := 0; i < len(blockpoolStrings); i++ {
		if(strings.Compare(blockpoolStrings[i], hash) == 0){
			if(i == len(blockpoolStrings) - 1){
				element = i + 1
			}
		}
	}

	fmt.Println(element)

	//return parent hash
	return blockpoolStrings[element]
}

func main() {
	
}
