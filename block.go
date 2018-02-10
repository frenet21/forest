package Block

import (
	"encoding/base64"
	"io"

	"golang.org/x/crypto/sha3"
)
import "crypto/cipher"
import "crypto/rsa"
import "crypto/aes"
import "crypto/rand"

type BlockData struct {
	encryptedKey     string   // RSA encrypted AES key
	encryptedMessage string   // AES encrypted message
	salt             [8]byte  // Random salt
	parent           [64]byte // Hash of parent block
}

type Block struct {
	// Block Data
	data BlockData

	// Non-hashed data
	ID     [64]byte // Hash of block data
	pepper [8]byte  // Random non-hashed salt
}

// Returns n random bytes
func RandomBytes(n int) []byte {
	out := make([]byte, n)
	return rand.Read(out)
}

// Selects a block parent based on the encrypted message
func selectParentHash(encryptedMessage string) string {
	// TODO: Connect this to the blockpool
	return base64.URLEncoding(sha3.New512().Sum(RandomBytes(32)))
}

func CreateBlockData(message string, *key rsa.PublicKey) BlockData {
	var out BlockData

	// Block salt
	out.salt = RandomBytes(8)

	// Message encryption
	// Random AES256 key
	AESkey := RandomBytes(32)
	// Block cipher for that key
	AESCipher := aes.NewCipher(AESkey)
	// encrypted message bytes
	cipherBytes := make([]bytes, aes.BlockSize+len(message))
	// Initialization Vector
	// Delivered with ciphertext as it is necessary for decryption...
	// But it doesn't have to be private to be secure
	iv := cipherBytes[:aes.BlockSize]
	// If error reading IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err) // panic!!
	}
	// Stream cipher
	stream := cipher.NewCTR(AESCipher, iv)
	// Plaintext bytes
	plaintext := []byte(message)
	// Encryption
	stream.XORKeyStream(cipherBytes[aes.BlockSize:], plaintext)
	// Convert to base64 and place in block
	out.encryptedMessage = base64.URLEncoding.EncodeToString(cipherBytes)

	// AES key encryption
	cipheredKey, e := rsa.EncryptOAEP(sha3.New512(), rand.Reader, key, AESkey, nil)
	// Panic on error
	if (e!=nil) {
		panic(err)
	}
	// Convert to base64 and place in block
	out.encryptedKey = base64.URLEncoding.EncodeToString(cipheredKey)
	
	// Select blockparent using blockpool
	out.parent=selectParentHash(out.encryptedMessage)
	
	// Done.
	return out
}

// BlockData -> string
func StringifyBlockData(data BlockData) string {
	// TODO implement this
	return ""
}

func DestringifyBlockData(data string) {
	var out BlockData
	
	// TODO Implement this

	return out
}

func CreateBlock (message string, *key rsa.PublicKey) Block {
	var out Block
	
	// Block data
	out.data = CreateBlockData(message, key)
	
	// Block ID
	dataString := StringifyBlockData(out.data)
	out.ID = base64.URLEncoding.EncodeToString(sha3.New512().Sum([]byte(dataString)))
	
	// Block pepper
	out.pepper = RandomBytes(8)
}
