package wordvalidator

import (
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/agnivade/levenshtein"
)

// IsSame returns if a is somewhat equal to b
func IsSame(a, b string) bool {
	if strings.EqualFold(a, b) {
		return true
	}

	maxDiff := 1

	minLen := min(len(a), len(b))
	if minLen > 48 {
		maxDiff = 8
	} else if minLen > 16 {
		maxDiff = 6
	} else if minLen > 7 {
		maxDiff = 4
	} else if minLen > 5 {
		maxDiff = 2
	}

	lenDiff := len(a) - len(b)
	if lenDiff > maxDiff || lenDiff < -maxDiff {
		return false
	}

	return levenshtein.ComputeDistance(a, b) <= maxDiff
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// NormalizeString normalizes the input.
// It changes the following things:
// - Space like characters are converted to spaces
// - Duplicated spaces are removed
// - Spaces around the string are removed
// - Non number and non letter characters are removed
// - Uppercase characters are converted to lowercase
func NormalizeString(inStr string) string {
	if len(inStr) == 0 {
		return ""
	}

	inBytes := []byte(inStr)

	for idx := len(inBytes) - 1; idx >= 0; idx-- {
		c := inBytes[idx]
		if c >= '0' && c <= '9' || c >= 'a' && c <= 'z' {
			continue
		}

		switch c {
		case '\n', '\r', '\t', ' ':
			if idx == len(inBytes)-1 {
				// Trim the space like characters around the input
				inBytes = inBytes[:idx]
			} else if idx != 0 && (inBytes[idx-1] == '\n' || inBytes[idx-1] == '\r' || inBytes[idx-1] == '\t' || inBytes[idx-1] == ' ') {
				// The character to the left is also a whitespace character, so we can remove this char
				// By doing this we remove the duplicated spaces
				inBytes = append(inBytes[:idx], inBytes[idx+1:]...)
			} else if c != ' ' {
				inBytes[idx] = ' '
			} else if idx == 0 {
				// The first character is a space, trim the front
				// We don't have to worry if the next character where spaces because
				// they would be already removed by a previous if else check
				inBytes = inBytes[1:]
			}
		default:
			if c >= 'A' && c <= 'Z' {
				// Convert uppercase to lowercase
				inBytes[idx] += 'a' - 'A'
			} else if c < utf8.RuneSelf {
				// Remove all other special ascii characters
				inBytes = append(inBytes[:idx], inBytes[idx+1:]...)
			}
		}
	}

	// Convert the inBytes to a string without copying the data
	return b2s(inBytes)
}

// b2s converts a byte slice to a string without copying
// Note that this will mean that changes made after to the byte slice will be reflected in the string and visa versa
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
