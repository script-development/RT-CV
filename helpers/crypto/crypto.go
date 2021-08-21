package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func NormalizeKey(key []byte) []byte {
	// Sum256 returns a bytes array of exactly 32 bytes, the length we also need for AES-256
	hashedKey := sha256.Sum256(key)
	return hashedKey[:]
}

func Encrypt(data, key []byte) ([]byte, error) {
	if len(key) < 16 {
		return nil, errors.New("an encryption keys needs to be at least 16 chars")
	}
	if data == nil {
		return nil, errors.New("data to encrypt should not be nil")
	}

	c, err := aes.NewCipher(NormalizeKey(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	encryptedValue, err := gcm.Seal(data[:0], nonce, data, nil), nil
	if err != nil {
		return nil, err
	}

	return append(nonce, encryptedValue...), nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) < 16 {
		return nil, errors.New("an decryption keys needs to be at least 16 chars")
	}

	c, err := aes.NewCipher(NormalizeKey(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
