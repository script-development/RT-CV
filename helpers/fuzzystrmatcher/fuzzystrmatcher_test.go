package fuzzystrmatcher

import (
	"os"
	"runtime/pprof"
	"testing"

	"github.com/apex/log"
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

		{"somewhere over the rainbow", "somewhere", false},
		{"banana", "apple", false},
	}

	for _, testCase := range testCases {
		parsedInput := Compile(false, testCase.input)
		if testCase.shouldMatch {
			True(t, parsedInput.Match(testCase.matchWith), `Expected "%s" to match "%v"`, testCase.input, testCase.matchWith)
		} else {
			False(t, parsedInput.Match(testCase.matchWith), `Expected "%s" to not match "%v"`, testCase.input, testCase.matchWith)
		}
	}

	matcher := Compile(false,
		"I love trees",
		"bananas are the best fruit",
		"banana",
	)

	matchesWith := []struct {
		input   string
		matches bool
	}{
		{"nothing", false},
		{"i love trees", true},
		{"bananas are the best fruit", true},
		{"banana", true},
		{"do you also love trees? i do.", true},
		{"on a sunday afternoon i like to eat a banana", true},
	}

	for _, testCase := range matchesWith {
		if testCase.matches {
			True(t, matcher.Match(testCase.input))
		} else {
			False(t, matcher.Match(testCase.input))
		}
	}
}

func BenchmarkMatch(b *testing.B) {
	f, err := os.Create("cpu.profile")
	if err != nil {
		log.WithField("error", err).Fatal("could not create cpu profile")
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.WithField("error", err).Fatal("could not start CPU profile")
	}

	defer func() {
		pprof.StopCPUProfile()
		_ = f.Close()
	}()

	matcher := Compile(false,
		"I love trees",
		"bananas are the best fruit",
		"banana",
	)

	matchesWith := []string{
		"nothing",
		"i love trees",
		"bananas are the best fruit",
		"banana",
		"do you also love trees? i do.",
		"on a sunday afternoon i like to eat a banana",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range matchesWith {
			matcher.Match(v)
		}
		matcher.Match(matchesWith...)
	}
}
