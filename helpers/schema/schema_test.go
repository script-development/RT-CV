package schema

import (
	"reflect"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestFromFailsWithWrongDataType(t *testing.T) {
	values := []interface{}{
		"foo",
		1,
		true,
	}

	for _, value := range values {
		_, _, err := From(reflect.TypeOf(value))
		Error(t, err)
	}
}

func TestFrom(t *testing.T) {
	type NestedStruct struct {
		B string
	}
	const constDummyArraySize = 32
	dummyArraySizeValue := uint(32)

	scenarios := []struct {
		name                   string
		in                     interface{}
		expectedProperties     map[string]Property
		expectedRequiredFields []string
	}{
		{
			"simple",
			struct{}{},
			map[string]Property{},
			[]string{},
		},
		{
			"with basic fields",
			struct {
				A string
				B int
				C bool
				D float64
			}{},
			map[string]Property{
				"A": {Type: PropertyTypeString},
				"B": {Type: PropertyTypeInteger},
				"C": {Type: PropertyTypeBoolean},
				"D": {Type: PropertyTypeNumber},
			},
			[]string{"A", "B", "C", "D"},
		},
		{
			"with json tag",
			struct {
				A string  `json:"renamed_field"`
				B float64 `json:"-"`
			}{},
			map[string]Property{
				"renamed_field": {Type: PropertyTypeString},
			},
			[]string{"renamed_field"},
		},
		{
			"with nested struct",
			struct {
				A NestedStruct
			}{},
			map[string]Property{
				"A": {
					Type: PropertyTypeObject,
					Properties: map[string]Property{
						"B": {Type: PropertyTypeString},
					},
					Required: []string{"B"},
				},
			},
			[]string{"A"},
		},
		{
			"with array and slice",
			struct {
				A []string
				B [constDummyArraySize]string
			}{},
			map[string]Property{
				"A": {Type: PropertyTypeArray},
				"B": {Type: PropertyTypeArray, MaxItems: &dummyArraySizeValue, MinItems: &dummyArraySizeValue},
			},
			[]string{},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			properties, requiredFields, err := From(reflect.TypeOf(s.in))
			NoError(t, err)
			Equal(t, s.expectedProperties, properties)
			Equal(t, s.expectedRequiredFields, requiredFields)
		})
	}
}
