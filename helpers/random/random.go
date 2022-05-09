package random

import (
	cryptoRand "crypto/rand"
	"math/rand"
)

const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars = "0123456789"
)

var allChars = lowerChars + upperChars + numberChars

// Seed the random number generator with a "good" random value
func Seed() {
	dst := make([]byte, 8)
	readBytes, err := cryptoRand.Read(dst)
	if err != nil {
		panic(err)
	}
	if readBytes != 8 {
		panic("did not read 8 random bytes")
	}

	seed := int64(0)
	for i, b := range dst {
		// Firstly we convert the 8 bit uint into a 8 bit int, so the bits stay in the first 8 bits
		// If we convert the uint8 directly to int64 the bits move into the 16bit range when the uint8 is more than 127
		bInt8 := int8(b)

		// Now we move the bits to the area inside the 64 bit int we want it to
		seed ^= (int64(bInt8) << (i * 8))
	}

	rand.Seed(seed)
}

// StringBytes returns a random set of string bytes
func StringBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = allChars[rand.Intn(len(allChars))]
	}
	return b
}

// SliceIndex returns a random slice index
func SliceIndex[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}

// GenerateKey generates a good random key
func GenerateKey() []byte {
	return StringBytes(32)
}
