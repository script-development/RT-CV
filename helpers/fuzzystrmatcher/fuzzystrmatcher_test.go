package fuzzystrmatcher

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestNormalizeString(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect string
	}{
		{"normalized inputs should not change", "abc", "abc"},
		{"should be converted to lowercase", "ABC", "abc"},
		{"spaces should be trimmed", "  ABC  ", "abc"},
		{"duplicated spaces should reduced to 1 space", "A  B   C", "a b c"},
		{"new line and tab characters should be replace by a space", "A\nB\tC", "a b c"},
		{"special characters should be removed", "a+b-c(d)", "abcd"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			Equal(t, testCase.expect, NormalizeString(testCase.input))
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		input       string
		matchWith   string
		shouldMatch bool
	}{
		{"a", "a", true},
		{"banana", "banana", true},
		{"banana", "banan", true},
		{"banana", "banaana", true},
		{"banana", "bananas", true},
		{"banana", "i want a banana", true},
		{"some thing", "thing some", true},
		{"says pet", "i love food says the pet", true},
		{"love i", "i love food says the pet", true},
		{
			"this is a very long sentence",
			"another sentence that contains the other sentence \"this is a very long sentence\" so there should be a match",
			true,
		},

		{"a", "b", false},
		{"some thing", "thing", false},
		{"banana", "apple", false},
	}

	for _, testCase := range testCases {
		parsedInput := Compile([]string{testCase.input})
		if testCase.shouldMatch {
			True(t, parsedInput.Match([]string{testCase.matchWith}), "Expected %s to match %v", testCase.input, testCase.matchWith)
		} else {
			False(t, parsedInput.Match([]string{testCase.matchWith}), "Expected %s to not match %v", testCase.input, testCase.matchWith)
		}
	}
}
