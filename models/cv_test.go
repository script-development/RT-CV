package models

import (
	"sort"
	"testing"
	"time"

	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/stretchr/testify/assert"
)

func TestSortEducations(t *testing.T) {
	now := time.Now()
	toArg := func(t time.Time) *jsonHelpers.RFC3339Nano {
		return jsonHelpers.RFC3339Nano(t).ToPtr()
	}

	tests := []struct {
		name     string
		in       Educations
		expected []string
	}{
		{
			"already sorted",
			Educations{
				{Name: "1", EndDate: toArg(now.AddDate(-1, 0, 0))},
				{Name: "2", EndDate: toArg(now.AddDate(-2, 0, 0))},
			},
			[]string{"1", "2"},
		},
		{
			"simple sort",
			Educations{
				{Name: "1", EndDate: toArg(now.AddDate(-2, 0, 0))},
				{Name: "2", EndDate: toArg(now.AddDate(-1, 0, 0))},
			},
			[]string{"2", "1"},
		},
		{
			"multiple",
			[]Education{
				{Name: "1", EndDate: toArg(now.AddDate(-2, 0, 0))},
				{Name: "2", EndDate: toArg(now.AddDate(-3, 0, 0))},
				{Name: "3", EndDate: toArg(now.AddDate(-1, 0, 0))},
			},
			[]string{"3", "1", "2"},
		},
		{
			"prefer endDate over startDate",
			[]Education{
				{Name: "1", EndDate: toArg(now.AddDate(-2, 0, 0))},
				{Name: "2", EndDate: toArg(now.AddDate(-3, 0, 0)), StartDate: toArg(now)},
				{Name: "3", EndDate: toArg(now.AddDate(-1, 0, 0))},
				{Name: "4", StartDate: toArg(now.AddDate(-2, -1, 0))},
			},
			[]string{"3", "1", "4", "2"},
		},
		{
			"place entries without a date at front and do not reorder them",
			[]Education{
				{Name: "1", EndDate: toArg(now.AddDate(-1, 0, 0))},
				{Name: "2"},
				{Name: "3", EndDate: toArg(now)},
				{Name: "4"},
			},
			[]string{"2", "4", "3", "1"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			sort.Sort(testCase.in)
			outNames := []string{}
			for _, edu := range testCase.in {
				outNames = append(outNames, edu.Name)
			}

			assert.Equal(t, testCase.expected, outNames)
		})
	}
}
