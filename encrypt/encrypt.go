package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	mrand "math/rand"
	"time"
)

const DefaultPublicKeyPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2OOeEUTCVa2PQfFA6fa1
kN4jg3fZtQDF6wY5aUW8/12JjJXQWFUssSB3fN9xXtZRSQYWnRDcz7S6c/i1Wbnd
4ArE4ulFv5QvEG+8PJsn0ruuBnBG6tSm3qjtCbN9lVR8xBzZbue9DYkwIKCGIk/q
2SywDYcoA/1moN2EV8YdQtsT480L92RzZ6hGSC2P9j4lZr/sjRGkXc4zudxKW6sx
eSWdo6RIPinitxnLlIduaDgA4xDZmH/X+2c2ged7oD1NdTx8o4arJjv1xBhcwuuQ
C1Lx+BF8qpDbLzAEeNN92zC9I9tyB4yNLGRVP9EF/3XxA4+0gmR+SDbCQb/MWfAw
xwIDAQAB
-----END PUBLIC KEY-----`
const DefaultPublickKeyID = 5097750102325959

// GenerateRangeNum return a random number between min and max
func GenerateRangeNum(min, max int) int {
	// mrand.Seed(time.Now().UnixNano())
	mrand.Seed(time.Now().Unix())
	randNum := mrand.Intn(max-min) + min
	return randNum
}

// GenerateRandomBytes return n bytes random data
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetPublicKey return the public key
func GetPublicKey(pubKeyPem []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pubKeyPem)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	return pubInterface.(*rsa.PublicKey)

}

func genEncryptKey(pubkey *rsa.PublicKey) (key, encKey []byte, err error) {
	key, _ = GenerateRandomBytes(16) // 128bits
	if err != nil {
		return
	}
	encKey, _ = rsa.EncryptPKCS1v15(rand.Reader, pubkey, key)

	return
}

// Encryptor is the excryptor struct
type Encryptor struct {
	pubKeyID     int64
	publicKey    *rsa.PublicKey
	symmetricKey []byte
	encKey       []byte
	version      int
	streamReader EncryptStreamReader
}

// EncryptStreamReader is the stream mode reader
type EncryptStreamReader interface {
	Read(p []byte) (n int, err error)
}

type cbcReader struct {
	r            io.Reader
	headerReader *bytes.Reader
	blockMode    cipher.BlockMode
	buf          *bytes.Buffer
	tmp          []byte
	err          error
}

// NewEncryptor return a Encryptor instance
func NewEncryptor(pubkey *rsa.PublicKey, version int) *Encryptor {
	if pubkey == nil {
		pubkey = GetPublicKey([]byte(DefaultPublicKeyPem))
	}

	aesKey, encKey, _ := genEncryptKey(pubkey)
	return &Encryptor{
		publicKey:    pubkey,
		symmetricKey: aesKey[:],
		encKey:       encKey[:],
		version:      version,
	}
}

func NewEncryptorV3(pubKeyID int64, pubkeyPEM []byte) *Encryptor {
	pubkey := GetPublicKey(pubkeyPEM)

	aesKey, encKey, _ := genEncryptKey(pubkey)
	return &Encryptor{
		pubKeyID:     pubKeyID,
		publicKey:    pubkey,
		symmetricKey: aesKey[:],
		encKey:       encKey[:],
		version:      3,
	}
}
