package auth

import (
	"encoding/base64"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAccessor can be used in tests to generate auth headers
// This is meant to be used only in tests
type TestAccessor struct {
	rollingKey string
	key        string
	salt       string
	keyAndSalt string // = keyBytes + saltBytes
	keyID      string
	authSeed   string
}

// NewAccessorHelper creates a TestAccessor that can be used to generate auth headers
func NewAccessorHelper(keyID primitive.ObjectID, key, salt, authSeed string) *TestAccessor {
	keyandSalt := key + salt

	h := hashSha(hashSha512, authSeed+keyandSalt)
	return &TestAccessor{
		rollingKey: h,
		key:        key,
		salt:       salt,
		keyAndSalt: keyandSalt,
		keyID:      keyID.Hex(),
		authSeed:   authSeed,
	}
}

// Key generates a new key
func (a *TestAccessor) Key() string {
	a.rollingKey = hashSha(hashSha512, a.rollingKey+a.keyAndSalt)

	token := strings.Join([]string{
		"sha512",
		a.keyID,
		a.salt,
		a.rollingKey,
	}, ":")

	return "Basic " + base64.URLEncoding.EncodeToString([]byte(token))
}
