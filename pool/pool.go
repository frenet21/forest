package pool

import (
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

type HashDate struct {
	hash string
	sent time.Time
}

type Blockpool struct {
	hashes [1000]string
	queue  []HashDate
}

var blockpool Blockpool // Blockpool singleton

// Returns the genesis pool, after setting the global pool to genesis
func genesisPool() Blockpool {
	var genesis Blockpool

	hasher := sha3.New512()

	for i := 0; i < len(genesis.hashes); i++ {
		block := make([]byte, 1)
		block[0] = byte(i)
		genesis.hashes[i] = string(hasher.Sum(block)[:64])
	}

	blockpool = genesis
	return genesis
}

// Adds a new hash to the blockpool's receive queue
func receiveBlockHash(hash string) {
	entry := HashDate{hash, time.Now()}
	blockpool.queue = append(blockpool.queue, entry)
}

func updateBlockpool() {
	firstOld := -1
	for i := 0; i < len(blockpool.queue); i++ {
		if time.Since(blockpool.queue[i].sent) > time.Hour {
			firstOld = i
			break
		}
	}
	if firstOld == -1 {
		return
	}
	pairs := blockpool.queue[firstOld:]
	hashes := make([]string, len(pairs))
	for i := 0; i < len(pairs); i++ {
		hashes[i] = pairs[i].hash
	}
	copy(blockpool.hashes[:], append(hashes, blockpool.hashes[:firstOld]...)[:1000])
}

// Selects a block parent based on the encrypted message
func selectParentHash(encryptedMessage string) string {
	updateBlockpool()

	//Add hash of encrypted message to the end of the blockpool array
	hash := string(sha3.New512().Sum([]byte(encryptedMessage))[:64])
	blockpoolStrings := append(blockpool.hashes[:], hash)

	//Sorted blockpool strings
	sort.Strings(blockpoolStrings)

	var element int

	//Determine the element after the encrypted message
	for i := 0; i < len(blockpoolStrings); i++ {
		if strings.Compare(blockpoolStrings[i], hash) == 0 {
			if i == len(blockpoolStrings)-1 {
				element = i + 1
			}
		}
	}

	//return parent hash
	return blockpoolStrings[element]
}

func main() {

}
