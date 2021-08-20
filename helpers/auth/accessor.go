package auth

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAccessor can be used in tests to generate auth headers
type TestAccessor struct {
	rollingKey      []byte
	keyBytes        []byte
	saltBytes       []byte
	keyAndSaltBytes []byte // = keyBytes + saltBytes
	keyId           string
	authSeed        []byte
}

func NewAccessorHelper(keyId primitive.ObjectID, key, salt string, authSeed []byte) *TestAccessor {
	keyBytes := []byte(key)
	saltBytes := []byte(salt)
	keyandSaltBytes := append(keyBytes, saltBytes...)

	h := sha512.Sum512(append(authSeed, keyandSaltBytes...))
	return &TestAccessor{
		rollingKey:      h[:],
		keyBytes:        keyBytes,
		saltBytes:       saltBytes,
		keyAndSaltBytes: keyandSaltBytes,
		keyId:           keyId.Hex(),
		authSeed:        authSeed,
	}
}

func (a *TestAccessor) Key() []byte {
	newRollingKey := sha512.Sum512(append(a.rollingKey, a.keyAndSaltBytes...))
	a.rollingKey = newRollingKey[:]

	src := bytes.Join([][]byte{
		[]byte("sha512:" + a.keyId),
		a.saltBytes,
		[]byte(hex.EncodeToString(a.rollingKey)),
	}, []byte(":"))

	return []byte("Basic " + base64.RawStdEncoding.EncodeToString(src))
}
