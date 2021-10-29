package crypto

import (
	"crypto/sha512"
	"encoding/hex"
)

// HashSha512String hashes a string using sha512 and returns the results hex encoded
func HashSha512String(in string) string {
	res := sha512.Sum512([]byte(in))
	return hex.EncodeToString(res[:])
}
