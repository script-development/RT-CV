package wordvalidator_test

import (
	"testing"

	"github.com/script-development/RT-CV/helpers/wordvalidator"
	"github.com/tj/assert"
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
			assert.Equal(t, testCase.expect, wordvalidator.NormalizeString(testCase.input))
		})
	}
}

type Profession struct {
	Name           string
	Alternatives   []string
	ForbiddenWords []string
}

var professions = []Profession{
	{
		Name: "verkoopmedewerker",
		Alternatives: []string{
			"verkoop medewerker",
			"VeRkOoP MeDeWeRkEr",
			"verkoop medewerkster",
			"verkoopmedewerkster",
			"verkoop   medewerker",
		},
		ForbiddenWords: []string{
			"medewerkster verkoop",
			"medewerker verkoop",
			"verkoop",
			"medewerker",
			"verkoper",
		},
	},
	{
		Name: "storingscoördinator",
		Alternatives: []string{
			"storingscoordinator",
			"storingscoördinator",
			"storings coordinator",
			"storings coördinator",
			"StoRinGs coördinator",
			"storing coordinator",
			"storing coördinator",
			"Storings coordinator",
		},
		ForbiddenWords: []string{
			"storings",
			"coordinator",
			"coördinator",
		},
	},
	{
		Name: "administratief medewerker",
		Alternatives: []string{
			"administratief medewerker",
			"administratief medewerkster",
			"administratieve medewerker",
		},
		ForbiddenWords: []string{
			"medewerker administratie",
			"medewerker administratieve",
			"administratie",
		},
	},
	{
		Name: "technisch administratief medewerker",
		Alternatives: []string{
			"technisch administratief medewerker",
			"technisch administratief medewerkster",
		},
		ForbiddenWords: []string{
			"administratief medewerker",
		},
	},
}

func TestCaseSentitivity(t *testing.T) {
	str1 := "AdmInIStrAtiEf medewerker"
	str2 := "administratief Medewerker"

	if !wordvalidator.IsSame(str1, str2) {
		t.Error("str1 and str2 aren't the same")
	}
}

var optimizationBypass bool

func BenchmarkIsSame(b *testing.B) {
	b.ReportAllocs()
	var same bool
	for i := 0; i < b.N; i++ {
		same = wordvalidator.IsSame("Administratief medewerkster", "Administratief medewerker")
	}
	optimizationBypass = same
}
