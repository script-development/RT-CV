package numbers

import "fmt"

// UintToBytes converst a uint64 into a number of bytes
// The size param defines the size of the output buffer
// The allowed sizes are:
// 	16: converts the in to uint16 and returns a 2 byte buffer
// 	32: converts the in to uint32 and returns a 4 byte buffer
// 	64: returns a 8 byte buffer
func UintToBytes(in uint64, size uint8) []byte {
	sizeInBytes := size / 8
	out := make([]byte, sizeInBytes)
	for i := uint8(0); i < sizeInBytes; i++ {
		out[i] = byte(in >> (i * 8))
	}
	return out
}

// BytesToUint convers a number of bytes into a uint64
// The allowed in lengths are:
// 	2: for a uint16 encoded as uint64
// 	4: for a uint32 encoded as uint64
// 	8: for a uint64
func BytesToUint(in []byte) (uint64, error) {
	inLen := uint8(len(in))
	switch inLen {
	case 2, 4, 8:
		// ok
	default:
		return 0, fmt.Errorf("Invalid byte length, expected 2, 4 or 8 bytes but got %d bytes", inLen)
	}

	out := uint64(0)
	for idx, b := range in {
		out |= (uint64(b) << (idx * 8))
	}

	return out, nil
}
