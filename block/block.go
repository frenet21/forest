package block

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"time"

	"github.com/stellar-tech/forest/pool"

	"golang.org/x/crypto/sha3"
)

type BlockData struct {
	encryptedKey     string   // RSA encrypted AES key
	encryptedMessage string   // AES encrypted message
	salt             [8]byte  // Random salt
	parent           [64]byte // Hash of parent block
	nonce            []byte   // Nonce used for GCM
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
	_, err := rand.Read(out)
	if err != nil {
		panic(err)
	}
	return out
}

// Selects a block parent based on the encrypted message
func selectParentHash(encryptedMessage string) [64]byte {
	// TODO: Connect this to the blockpool
	var out [64]byte
	copy(out[:], sha3.New512().Sum(RandomBytes(32))[:64])
	return out
}

func CreateBlockData(message string, key *rsa.PublicKey) BlockData {
	var out BlockData

	// Controls how long we wait for encryption to complete
	// Go doesn't perform encryptions in constant-time...
	// So to prevent timing attacks, we wait after encryption
	// The time is taken before running the encryption, and then after encrypt
	//    we wait until that much time has elapsed
	// Thus, we get pseudo-constant time behavior
	// This time needs to be long enough that encryption of the key and of the
	//    message will be complete, each in one period, for any (reasonable) message.
	constantDelayFactor := 500 * time.Millisecond

	// Block salt
	copy(out.salt[:], RandomBytes(8)[:8])

	// Message encryption
	// Random AES256 key
	AESkey := RandomBytes(32)
	// Block cipher for that key
	AESCipher, e := aes.NewCipher(AESkey)
	if e != nil {
		panic(e)
	}
	// AEAD
	auth, err := cipher.NewGCM(AESCipher)
	if err != nil {
		panic(err)
	}
	// Initialization Vector
	// Delivered with ciphertext as it is necessary for decryption...
	// But it doesn't have to be private to be secure
	out.nonce = RandomBytes(auth.NonceSize())
	// Plaintext bytes
	plaintext := []byte(message)

	// Encryption
	// First, get current time and add delay factor
	endpoint := time.Now().Add(constantDelayFactor)
	// Then actually run encryption
	cipherBytes := auth.Seal(nil, out.nonce, plaintext, out.salt[:])
	// Now, delay until we reach endpoint
	time.Sleep(time.Until(endpoint))

	// Convert to base64 and place in block
	out.encryptedMessage = base64.URLEncoding.EncodeToString(cipherBytes)

	// Select blockparent using blockpool
	out.parent = selectParentHash(out.encryptedMessage)

	// AES key encryption
	// First, get current time and add delay factor
	endpoint = time.Now().Add(constantDelayFactor)
	// Then actually run encryption
	cipheredKey, e := rsa.EncryptOAEP(sha3.New512(), rand.Reader, key, AESkey, out.parent[:])
	// Now, delay until we reach endpoint
	time.Sleep(time.Until(endpoint))

	// Panic on error
	if e != nil {
		panic(e)
	}
	// Convert to base64 and place in block
	out.encryptedKey = base64.URLEncoding.EncodeToString(cipheredKey)

	// Done.
	return out
}

// BlockData -> string
func StringifyBlockData(data BlockData) string {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(data)
	raw := buf.Bytes()
	return string(raw[:buf.Len()])
}

func DestringifyBlockData(data string) BlockData {
	var out BlockData

	var buf bytes.Buffer
	buf.WriteString(data)
	decoder := gob.NewDecoder(&buf)
	decoder.Decode(out)

	return out
}

func CreateBlock(message string, key *rsa.PublicKey) Block {
	var out Block

	// Block data
	// This is where encryption is done...
	// Constant factor delay?
	out.data = CreateBlockData(message, key)

	// Block ID
	dataString := StringifyBlockData(out.data)
	copy(out.ID[:], sha3.New512().Sum([]byte(dataString))[:64])

	// Block pepper
	copy(out.pepper[:], RandomBytes(8)[:8])

	return out
}

func StringifyBlock(block Block) string {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(block)
	raw := buf.Bytes()
	return string(raw[:buf.Len()])
}

func DestringifyBlock(block string) Block {
	var out Block

	var buf bytes.Buffer
	buf.WriteString(block)
	decoder := gob.NewDecoder(&buf)
	decoder.Decode(out)

	return out
}

// Returns the decrypted message from a block with a given PrivateKey
func AttemptDecrypt(block Block, key *rsa.PrivateKey) (message string, err error) {
	// First off, we confirm the integrity of the block data
	// If the blockID doesn't match the hash of the blockdata, then it has been modified
	// If that occurs, report an error
	dataString := StringifyBlockData(block.data)
	test := sha3.New512().Sum([]byte(dataString))
	if !bytes.Equal(test, block.ID[:]) {
		// The blockdata has been modified!
		// Error out
		return "", errors.New("Blockdata hash mismatch: ID " + string(block.ID[:64]) +
			" is not equal to hash of data " + string(test[:64]))
	}
	// No data tampering has occurred if we get here...
	// Or if it has, it caused a collision in SHA3-512, which is insanely unlikely

	// Controls how long we wait for decryption to complete
	// Go doesn't perform encryptions in constant-time...
	// So to prevent timing attacks, we wait after decryption
	// The time is taken before running the decryption, and then after encrypt
	//    we wait until that much time has elapsed
	// Thus, we get pseudo-constant time behavior
	// This time needs to be long enough that decryption of the key and of the
	//    message will be complete, each in one period, for any (reasonable) message.
	constantDelayFactor := 500 * time.Millisecond

	// First, attempt to decrypt the encryptedKey
	// First, get our encrypted key as a byte array
	encryptedKeyBytes, er := base64.URLEncoding.DecodeString(block.data.encryptedKey)
	if er != nil {
		return "", er
	}
	// Then, get current time and add constant delay factor
	endpoint := time.Now().Add(constantDelayFactor)
	// Then actually attempt decryption
	AESkey, e := rsa.DecryptOAEP(sha3.New512(), rand.Reader, key, encryptedKeyBytes, block.data.parent[:])
	// Now wait until endpoint
	time.Sleep(time.Until(endpoint))
	// Return on error
	if e != nil {
		return "", e
	}

	// Now, attempt to use that key to decrypt the encryptedMessage
	AESCipher, err := aes.NewCipher(AESkey)
	if err != nil {
		return "", err
	}
	msg, error := base64.URLEncoding.DecodeString(block.data.encryptedMessage)
	if error != nil {
		return "", error
	}
	auth, er := cipher.NewGCM(AESCipher)
	if er != nil {
		return "", er
	}
	// First, get current time and add constant delay factor
	endpoint = time.Now().Add(constantDelayFactor)
	// Now actually attempt decryption
	msg, error = auth.Open(nil, block.data.nonce, msg, block.data.salt[:])
	// Now wait until endpoint
	time.Sleep(time.Until(endpoint))
	if error != nil {
		return "", error
	}
	// And return
	return string(msg), nil
}

// Call on main startup
func Initialize() {
	gob.Register(BlockData{})
	gob.Register(Block{})
}
