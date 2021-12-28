package crypto

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	. "github.com/stretchr/testify/assert"
)

func TestNonceChain(t *testing.T) {
	nonce := []byte("test")
	hasher := sha512.New()
	genNextNonce(hasher, nonce)
	Equal(
		t,
		"ee26b0dd",
		hex.EncodeToString(nonce),
		"nonce should be hashed and the result hash should have the length of the input nonce",
	)
	genNextNonce(hasher, nonce)
	Equal(t, "6c23a22c", hex.EncodeToString(nonce), "nonce should be diffrent after genNextNonce")
	genNextNonce(hasher, nonce)
	Equal(t, "78e9ae46", hex.EncodeToString(nonce), "nonce should be diffrent after genNextNonce")

	nonce = []byte("a diffrent test")
	genNextNonce(hasher, nonce)
	Equal(t, "3f6ca625423a948fa4479f1bab5af1", hex.EncodeToString(nonce), "expect a total diffrent nonce after giving it another input")
	genNextNonce(hasher, nonce)
	Equal(t, "d46b6793ad029f62519e9647aa20b4", hex.EncodeToString(nonce))
	genNextNonce(hasher, nonce)
	Equal(t, "bed3069ee86b00a338aec0ef1237a2", hex.EncodeToString(nonce))
}

func TestEncryptWriterDecryptReader(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"simple data", "very secret data that should be encrypted"},
		{"large data", strings.Repeat("very nice.", chunkSize)},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			// ENCRYPT
			key := []byte("a-testkey-that-is-longer-than-16-chars")
			resultBuffer := bytes.NewBuffer(nil)
			encryptWriter, err := NewEncryptWriter(key, resultBuffer)
			NoError(t, err)
			NotNil(t, encryptWriter)
			lenOfBuffBeforeWriting := resultBuffer.Len()

			n, err := encryptWriter.Write([]byte(testCase.data))
			NoError(t, err)
			Equal(t, len(testCase.data), n)

			err = encryptWriter.Close()
			NoError(t, err)

			NotEqual(t, 0, resultBuffer.Len())
			NotEqual(t, lenOfBuffBeforeWriting, resultBuffer.Len())

			// DECRYPT
			encryptReader, err := NewEncryptReader(key, resultBuffer)
			NoError(t, err)

			decryptedData := make([]byte, len(testCase.data))
			n, err = io.ReadFull(encryptReader, decryptedData)
			NoError(t, err)
			Equal(t, len(testCase.data), n)
			Equal(t, len(testCase.data), len(decryptedData))
			True(t, utf8.Valid(decryptedData))
			if testCase.data != string(decryptedData) {
				panic("decrypted data is not equal to the original data")
			}
		})
	}
}
