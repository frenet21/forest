package pool

import (
	//"golang.org/x/crypto/sha3"
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
	//Added encrypted message to the end of the blockpool array
	blockpoolStrings := append(info.blockpool, info.encryptedMessage)

	//Sorted blockpool strings
	fmt.Println(blockpoolStrings)
	sort.Strings(blockpoolStrings)
	fmt.Println(blockpoolStrings)

	var element int

	//Determine the element after the encrypted message
	for i := 0; i < len(blockpoolStrings); i++ {
		if(strings.Compare(blockpoolStrings[i], info.encryptedMessage) == 0){
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
