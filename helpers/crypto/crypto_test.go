package crypto

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	data := "very secret data that should be encrypted"
	key := []byte("a-testkey-that-is-longer-than-16-chars")
	encryptedData, err := Encrypt([]byte(data), key)
	NoError(t, err)
	NotNil(t, encryptedData)
	NotEqual(t, []byte(data), encryptedData)

	decryptedData, err := Decrypt(encryptedData, key)
	NoError(t, err)
	Equal(t, data, string(decryptedData))

	decryptedData, err = Decrypt(encryptedData, []byte("this is an invalid key and should fail"))
	Error(t, err)
	NotEqual(t, data, string(decryptedData))
}
