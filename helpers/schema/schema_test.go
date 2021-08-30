package schema

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestFromFailsWithWrongDataType(t *testing.T) {
	values := []interface{}{
		"foo",
		1,
		true,
		nil,
	}

	for _, value := range values {
		_, err := From(
			value,
			"",
			func(string, Property) {},
			func(string) bool { return true },
			&WithMeta{},
		)
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
			"with jsonSchema tag",
			struct {
				A *string
				B *string `jsonSchema:"required"`
				C string
				D string   `jsonSchema:"notRequired"`
				E string   `jsonSchema:"notRequired,deprecated"`
				F []string `jsonSchema:"uniqueItems"`
			}{},
			map[string]Property{
				"A": {Type: PropertyTypeString},
				"B": {Type: PropertyTypeString},
				"C": {Type: PropertyTypeString},
				"D": {Type: PropertyTypeString},
				"E": {Type: PropertyTypeString, Deprecated: true},
				"F": {Type: PropertyTypeArray, Items: &Property{Type: PropertyTypeString}, UniqueItems: true},
			},
			[]string{"B", "C"},
		},
		{
			"with nested struct",
			struct {
				A NestedStruct
			}{},
			map[string]Property{
				"A": {
					Ref: "#/testing/HelpersSchemaNestedStruct",
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
				"A": {
					Type: PropertyTypeArray,
					Items: &Property{
						Type: PropertyTypeString,
					},
				},
				"B": {
					Type:     PropertyTypeArray,
					MaxItems: &dummyArraySizeValue,
					MinItems: &dummyArraySizeValue,
					Items: &Property{
						Type: PropertyTypeString,
					},
				},
			},
			[]string{},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			property, err := From(
				s.in,
				"#/testing/",
				func(key string, property Property) {},
				func(key string) bool { return true },
				nil,
			)
			NoError(t, err)
			Equal(t, s.expectedProperties, property.Properties)
			Equal(t, s.expectedRequiredFields, property.Required)
		})
	}
}
