package block

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

type block struct {
	// Block Data
	data BlockData

	// Non-hashed data
	ID     [64]byte // Hash of block data
	pepper [8]byte  // Random non-hashed salt
}

// Returns n random bytes
func RandomBytes(n int) ([]byte, error) {
	out := make([]byte, n)
	return rand.Read(out)
}

func (message string, parent [64]byte, key rsa.PublicKey) CreateBlockData() BlockData {
	out := BlockData

	// Block salt
	out.salt = RandomBytes(8)

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

}
