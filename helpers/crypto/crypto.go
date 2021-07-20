package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"time"
)

type Encrypter struct {
	// This can be used to identify if data is encrypted using this encryption instance
	// This is a random value separate from the nonce and chiper so this can only be used to identify the encryption instance
	MasterKeyID []byte

	chiper cipher.AEAD

	// A random sha256 hash that is hashed by itself every encryption
	//
	// Why hasing it by itself?
	// rand.Reader is often very slow so we hash this value by itself to get a pseudo-random string back
	//
	// Is the above secure?
	// Reading https://en.wikipedia.org/wiki/Cryptographic_nonce#Definition we are allowed to use pseudo-random nonce values
	fullNonce [32]byte
}

func Init() (Encrypter, error) {
	res := Encrypter{}

	// generate master key id
	fullMashterKeyId := sha256.Sum256([]byte(time.Now().Format(time.RFC3339Nano)))
	res.MasterKeyID = fullMashterKeyId[:16]

	// generate nonce
	nonceSeed := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, nonceSeed)
	if err != nil {
		return res, err
	}
	res.fullNonce = sha256.Sum256(nonceSeed)

	// Generate key
	key := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, key)
	if err != nil {
		return res, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return res, err
	}
	res.chiper, err = cipher.NewGCM(c)

	return res, err
}

// Encrypt encrypts the input data
// Returns:
// [16 KeyID] + [12 nonce] + [... encrypted data]
func (e *Encrypter) Encrypt(data []byte) []byte {
	e.fullNonce = sha256.Sum256(e.fullNonce[:])
	return append(
		e.MasterKeyID,
		append(
			e.fullNonce[:12],
			e.chiper.Seal(nil, e.fullNonce[:12], data, nil)...,
		)...,
	)
}

func (e *Encrypter) Decrypt(data []byte) ([]byte, error) {
	masterKeyLen := len(e.MasterKeyID)
	if len(data) < masterKeyLen+12+1 {
		return nil, errors.New("invalid data length")
	}

	if !bytes.Equal(data[0:masterKeyLen], e.MasterKeyID) {
		return nil, errors.New("encrypted using another key")
	}

	return e.chiper.Open(nil, data[masterKeyLen:masterKeyLen+12], data[masterKeyLen+12:], nil)
}
