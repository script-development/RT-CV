package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimizeQuery(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"", ""},
		{" ", ""},
		{"a", "a"},
		{"foo", "foo"},
		{"foo bar", "foo bar"},
		{"foo\tbar", "foo bar"},
		{"FOO BAR", "foo bar"},
	}

	for _, tCase := range testCases {
		out, _ := optimizeQuery(tCase.in)
		assert.Equal(t, tCase.out, out)
	}
}
