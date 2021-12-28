package numbers

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestUintToByteToUint(t *testing.T) {
	testCases := []struct {
		value            uint64
		isMoreThan32Bits bool
	}{
		{0, false},
		{1, false},
		{500, false},
		{123456789, false},
		{1234567890123456789, true},
	}

	for _, testCase := range testCases {
		bytes := UintToBytes(testCase.value, 64)
		resultValue, err := BytesToUint(bytes[:])
		NoError(t, err)
		Equal(t, testCase.value, resultValue)

		bytes = UintToBytes(testCase.value, 32)
		resultValue, err = BytesToUint(bytes[:])
		NoError(t, err)
		if testCase.isMoreThan32Bits {
			NotEqual(t, testCase.value, resultValue)
		} else {
			Equal(t, testCase.value, resultValue)
		}
	}
}
