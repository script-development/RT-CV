package validation

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestValidDomain(t *testing.T) {
	cases := []struct {
		valid         bool
		domain        string
		allowWildcard bool
	}{
		{true, "example.com", false},
		{true, "a-a.com", false},
		{true, "a.a.com", false},
		{true, "*.example.com", true},
		{true, "test.*.example.com", true},
		{true, "*", true},
		{false, "*", false},
		{false, "example", false},
		{false, "", false},
		{false, ".", false},
		{false, ".com", false},
		{false, ".example.com", false},
		{false, "example..com", false},
		{false, "-.com", false},
		{false, "example.-", false},
		{false, "example.-.com", false},
		{false, "example.-.com", false},
		{false, "example.-.com", false},
		{false, "-example.com", false},
		{false, "example-.com", false},
	}
	for _, testcase := range cases {
		err := ValidDomain(testcase.domain, testcase.allowWildcard)
		if testcase.valid {
			NoError(t, err, testcase.domain)
		} else {
			Error(t, err, testcase.domain)
		}
	}
}
