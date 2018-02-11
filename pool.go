package main

import (
	"bytes"
	"encoding/gob"
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
	Hashes [1000]string
	queue  []HashDate
}

var blockpool Blockpool // Blockpool singleton

// Returns the genesis pool, after setting the global pool to genesis
func GenesisPool() Blockpool {
	var genesis Blockpool

	hasher := sha3.New512()

	for i := 0; i < len(genesis.Hashes); i++ {
		block := make([]byte, 1)
		block[0] = byte(i)
		hasher.Write(block)
		genesis.Hashes[i] = string(hasher.Sum(nil)[:64])
		hasher.Reset()
	}

	blockpool = genesis
	return genesis
}

// Adds a new hash to the blockpool's receive queue
func ReceiveBlockHash(hash string) {
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
	copy(blockpool.Hashes[:], append(hashes, blockpool.Hashes[:firstOld]...)[:1000])
}

// Selects a block parent based on the encrypted message
func SelectParentHash(encryptedMessage string) string {
	updateBlockpool()

	//Add hash of encrypted message to the end of the blockpool array
	hasher := sha3.New512()
	hasher.Write([]byte(encryptedMessage))
	hash := string(hasher.Sum(nil)[:64])
	blockpoolStrings := append(blockpool.Hashes[:], hash)

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

func StringifyBlockpool() string {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(blockpool)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Returns the blockpool after assigning it
func DestringifyBlockpool(pool string) Blockpool {
	var buf bytes.Buffer
	buf.WriteString(pool)
	decoder := gob.NewDecoder(&buf)
	err := decoder.Decode(blockpool)
	if err != nil {
		panic(err)
	}

	return blockpool
}
