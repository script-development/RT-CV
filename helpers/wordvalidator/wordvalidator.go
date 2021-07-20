package wordvalidator

import (
	"strings"

	"github.com/agnivade/levenshtein"
)

func IsSame(a, b string) bool {
	return strings.EqualFold(a, b) || levenshtein.ComputeDistance(strings.TrimSpace(a), strings.TrimSpace(b)) <= 1 // 1 allocation :(
}
