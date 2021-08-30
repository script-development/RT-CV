package schema

import (
	"errors"
	"reflect"
	"strings"
)

// Describe can be implmented by a type to manually describe the type
type Describe interface {
	JSONSchemaDescribe() Property
}

var reflectDescribe = reflect.TypeOf((*Describe)(nil)).Elem()

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
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Type        PropertyType        `json:"type,omitempty"`  // The data type
	Enum        []interface{}       `json:"enum,omitempty"`  // The value should validate againest one of these
	Const       interface{}         `json:"const,omitempty"` // Equal to a enum with 1 value
	Deprecated  bool                `json:"deprecated,omitempty"`
	Default     interface{}         `json:"default,omitempty"`
	Examples    []interface{}       `json:"examples,omitempty"`
	Format      Format              `json:"format,omitempty"`
	Ref         string              `json:"$ref,omitempty"`
	Defs        map[string]Property `json:"$defs,omitempty"`

	// Only in the root of the schema
	Schema Version `json:"$schema,omitempty"`
	ID     string  `json:"$id,omitempty"`

	// type == object
	Properties        map[string]Property `json:"properties,omitempty"` // required field
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

var errInvalidSchemaFromInput = errors.New("argument must be a struct, map, slice or array or a pointer to one of those")

// WithMeta adds the Schema and ID field to the returned property
type WithMeta struct {
	SchemaID string
}

// From converts a struct into a value for the  properties part of a json schema
// baseRefPath might look something like #/components/schemas/
func From(
	inputType interface{},
	baseRefPath string,
	addRef func(key string, property Property),
	hasRef func(key string) bool,
	meta *WithMeta,
) (Property, error) {
	if inputType == nil {
		return Property{}, errInvalidSchemaFromInput
	}

	t := reflect.TypeOf(inputType)
outerLoop:
	for {
		switch t.Kind() {
		case reflect.Ptr:
			t = t.Elem()
		default:
			break outerLoop
		}
	}

	var res Property
	switch t.Kind() {
	case reflect.Struct:
		properties, requiredFields := parseStruct(t, baseRefPath, addRef, hasRef)
		res = Property{
			Type:       PropertyTypeObject,
			Properties: properties,
			Required:   requiredFields,
		}
	case reflect.Map, reflect.Array, reflect.Slice:
		res, _, _ = parseType(t, baseRefPath, addRef, hasRef)
	default:
		return Property{}, errInvalidSchemaFromInput
	}
	if meta != nil {
		res.Schema = VersionUsed
		res.ID = meta.SchemaID
	}
	return res, nil
}

func parseStruct(
	t reflect.Type,
	baseRefPath string,
	addRef func(key string, property Property),
	hasRef func(key string) bool,
) (properties map[string]Property, requiredFields []string) {
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

		argRequired := false
		argNotRequired := false
		argDeprecated := false
		argUniqueItems := false
		args := strings.Split(field.Tag.Get("jsonSchema"), ",")
		for _, arg := range args {
			switch arg {
			case "notRequired":
				argNotRequired = true
			case "required":
				argRequired = true
			case "deprecated":
				argDeprecated = true
			case "uniqueItems":
				argUniqueItems = true
			}
		}

		property, required, skip := parseType(field.Type, baseRefPath, addRef, hasRef)
		if skip {
			continue
		}

		if argRequired {
			required = true
		}
		if argNotRequired {
			required = false
		}
		if argDeprecated {
			property.Deprecated = true
		}
		if argUniqueItems && property.Type == PropertyTypeArray {
			property.UniqueItems = true
		}

		properties[name] = property
		if required {
			requiredFields = append(requiredFields, name)
		}
	}

	return properties, requiredFields
}

func parseType(
	t reflect.Type,
	baseRefPath string,
	addRef func(key string, property Property),
	hasRef func(key string) bool,
) (property Property, required bool, skip bool) {
	required = true

	for {
		if t.Kind() != reflect.Ptr {
			break
		}
		required = false
		t = t.Elem()
	}

	if t.Implements(reflectDescribe) {
		methodName := "JSONSchemaDescribe"
		valueOfT := reflect.New(t).Elem()
		output := valueOfT.MethodByName(methodName).Call([]reflect.Value{})[0]
		property, ok := output.Interface().(Property)
		if !ok {
			panic("method " + methodName + " did not return the expected value type")
		}
		return property, required, false
	}

	switch t.Kind() {
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
		arrayLen := uint(t.Len())
		property.MinItems = &arrayLen
		property.MaxItems = &arrayLen
		fallthrough
	case reflect.Slice:
		property.Type = PropertyTypeArray
		required = false
		innerProperty, _, skip := parseType(t.Elem(), baseRefPath, addRef, hasRef)
		if skip {
			return property, required, true
		}
		property.Items = &innerProperty
	case reflect.Interface:
	case reflect.Map:
		// TODO maybe there is some way have not strictly defined keys but strictly defined values
		property.Type = PropertyTypeObject
	case reflect.Ptr:
	case reflect.String:
		property.Type = PropertyTypeString
	case reflect.Struct:
		key := ""
		if t.Name() != "" {
			parts := append(strings.Split(t.PkgPath(), "/")[3:], t.Name())
			for idx, part := range parts {
				// convert every part first letter to an uppercase
				parts[idx] = strings.ToUpper(part[0:1]) + part[1:]
			}
			key = strings.Join(parts, "")
		}

		if key == "" || !hasRef(key) {
			// This is here for when a same type struct is embedded inside of itself
			addRef(key, property)

			properties, required := parseStruct(t, baseRefPath, addRef, hasRef)
			property = Property{
				Required:   required,
				Type:       PropertyTypeObject,
				Properties: properties,
			}
			if key != "" {
				addRef(key, property)
				property = Property{Ref: baseRefPath + key}
			}
		} else {
			property = Property{Ref: baseRefPath + key}
		}
	case reflect.Chan, reflect.Func:
		// These fields are ignored by json marshall so we do to
		return property, required, true
	}

	return property, required, false
}
