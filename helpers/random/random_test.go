package random

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestSeedShouldNotPanic(*testing.T) {
	Seed()
}

func TestStringBytes(t *testing.T) {
	randomValue := StringBytes(1)
	NotNil(t, randomValue)
	Len(t, randomValue, 1)

	randomValue = StringBytes(64)
	NotNil(t, randomValue)
	Len(t, randomValue, 64)

	otherRandomValue := StringBytes(64)
	NotEqual(t, randomValue, otherRandomValue)
}

func TestGenerateKey(t *testing.T) {
	key := GenerateKey()
	NotNil(t, key)
	Len(t, key, 32)
}
