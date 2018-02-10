package block

/* import "golang.org/x/crypto/sha3"
import "crypto/cipher"
import "crypto/rsa"
import "crypto/aes"
import "crypto/rand" */

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
