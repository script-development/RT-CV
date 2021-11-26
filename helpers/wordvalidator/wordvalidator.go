package wordvalidator

import (
	"strings"

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
