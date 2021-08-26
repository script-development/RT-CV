package schema

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// PropertyType contains the value type of a property
type PropertyType string

const (
	// PropertyTypeNull represends nil values
	PropertyTypeNull = PropertyType("null")
	// PropertyTypeBoolean represends boolean values
	PropertyTypeBoolean = PropertyType("boolean")
	// PropertyTypeObject represends struct and map values
	PropertyTypeObject = PropertyType("object")
	// PropertyTypeArray represends slice and array values
	PropertyTypeArray = PropertyType("array")
	// PropertyTypeInteger represends int and uint values
	PropertyTypeInteger = PropertyType("integer")
	// PropertyTypeNumber represends float values
	PropertyTypeNumber = PropertyType("number")
	// PropertyTypeString represends string values
	PropertyTypeString = PropertyType("string")
)

// Format contains a expected data format
type Format string

const (
	// FormatDateTime derives from RFC 3339
	FormatDateTime = Format("date-time")
	// FormatDate derives from RFC 3339
	FormatDate = Format("date")
	// FormatTime derives from RFC 3339
	FormatTime = Format("time")
	// FormatDuration derives from RFC 3339
	FormatDuration = Format("duration")
	// FormatEmail as defined by the "Mailbox" ABNF rule in RFC 5321, section 4.1.2.
	FormatEmail = Format("email")
	// FormatIdnEmail as defined by the extended "Mailbox" ABNF rule in RFC 6531, section 3.3.
	FormatIdnEmail = Format("idn-email")
	// FormatHostname as defined by RFC 1123, section 2.1, including host names produced using the Punycode
	// algorithm specified in RFC 5891, section 4.4.
	FormatHostname = Format("hostname")
	// FormatIdnHostname as defined by either RFC 1123 as for hostname, or an internationalized hostname as defined
	// by RFC 5890, section 2.3.2.3.
	FormatIdnHostname = Format("idn-hostname")
	// FormatIPV4 is a ip version 4
	FormatIPV4 = Format("ipv4")
	// FormatIPV6 is a ip version 6
	FormatIPV6 = Format("ipv6")
	// FormatURI a valid uri derived from RFC3986
	FormatURI = Format("uri")
	// FormatURIReference a valid uri derived from RFC3986
	// either a URI or a relative-reference
	FormatURIReference = Format("uri-reference")
	// FormatIRI a valid uri derived from RFC3987
	FormatIRI = Format("iri")
	// FormatIRIReference a valid uri derived from RFC3987
	// either an IRI or a relative-reference
	FormatIRIReference = Format("iri-reference")
	// FormatUUID a valid uuid derived from RFC4122
	FormatUUID = Format("uui")
)

// Version defines the version of the schema
type Version string

// VersionUsed contains the schema version this package was build ontop
const VersionUsed = Version("https://json-schema.org/draft/2020-12/schema")

// Property represends a map / struct entry
type Property struct {
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Type        PropertyType  `json:"type,omitempty"`  // The data type
	Enum        []interface{} `json:"enum,omitempty"`  // The value should validate againest one of these
	Const       interface{}   `json:"const,omitempty"` // Equal to a enum with 1 value
	Deprecated  bool          `json:"deprecated"`
	Default     interface{}   `json:"default,omitempty"`
	Examples    []interface{} `json:"examples,omitempty"`
	Format      Format        `json:"format,omitempty"`
	Schema      Version       `json:"$schema"`
	ID          string        `json:"$id"`

	// type == object
	Properties        map[string]Property // required field
	Required          []string            `json:"required,omitempty"`
	MaxProperties     *uint               `json:"maxProperties,omitempty"`
	MinProperties     *uint               `json:"minProperties,omitempty"`
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"`

	// type == number || type == integer
	Minimum          *int `json:"minimum,omitempty"`          // >=
	Maximum          *int `json:"maximum,omitempty"`          // <=
	ExclusiveMinimum *int `json:"exclusiveMinimum,omitempty"` // >
	ExclusiveMaximum *int `json:"exclusiveMaximum,omitempty"` // <
	MultipleOf       uint `json:"multipleOf,omitempty"`

	// type == array
	Items       *Property `json:"items,omitempty"` // required field
	MinItems    *uint     `json:"minItems,omitempty"`
	MaxItems    *uint     `json:"maxItems,omitempty"`
	UniqueItems bool      `json:"uniqueItems,omitempty"`
	MaxContains *uint     `json:"maxContains,omitempty"`
	MinContains *uint     `json:"minContains,omitempty"`

	// type == string
	MaxLength *uint  `json:"maxLength,omitempty"`
	MinLength *uint  `json:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"` // ECMA-262 regular expression

}

var errInvalidSchemaFromInput = errors.New("argument must be a struct or a poitner to a struct")

// From converts a struct into a value for the  properties part of a json schema
// FIXME add cache to already converted types
func From(t reflect.Type) (properties map[string]Property, requiredFields []string, err error) {
	for {
		if t.Kind() != reflect.Ptr {
			break
		}
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, nil, errInvalidSchemaFromInput
	}

	requiredFields = []string{}
	properties = map[string]Property{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		jsonTagParts := strings.Split(jsonTag, ",")
		name := field.Name
		customName := jsonTagParts[0]
		if len(customName) > 0 {
			if customName == "-" {
				continue
			}
			name = customName
		}

		required := true
		property := Property{}

		fieldType := field.Type
		for {
			if fieldType.Kind() != reflect.Ptr {
				break
			}
			required = false
			fieldType = fieldType.Elem()
		}

		switch fieldType.Kind() {
		case reflect.Bool:
			property.Type = PropertyTypeBoolean
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fallthrough
		case reflect.Complex64, reflect.Complex128:
			fallthrough
		case reflect.UnsafePointer:
			property.Type = PropertyTypeInteger
		case reflect.Float32, reflect.Float64:
			property.Type = PropertyTypeNumber
		case reflect.Array:
			arrayLen := uint(fieldType.Len())
			property.MinItems = &arrayLen
			property.MaxItems = &arrayLen
			fallthrough
		case reflect.Slice:
			property.Type = PropertyTypeArray
			required = false
		case reflect.Interface:
		case reflect.Map:
			// TODO maybe there is some way have not strictly defined keys but strictly defined values
			property.Type = PropertyTypeObject
		case reflect.Ptr:
		case reflect.String:
			property.Type = PropertyTypeString
		case reflect.Struct:
			property.Type = PropertyTypeObject
			properties, requiredFields, err := From(fieldType)
			if err != nil {
				return nil, nil, fmt.Errorf("field %s failed with error: %s", field.Name, err.Error())
			}
			property.Properties = properties
			property.Required = requiredFields
		case reflect.Chan, reflect.Func:
			// These fields are ignored by json marshall so we do to
			continue
		}

		properties[name] = property
		if required {
			requiredFields = append(requiredFields, name)
		}
	}

	return properties, requiredFields, nil
}
