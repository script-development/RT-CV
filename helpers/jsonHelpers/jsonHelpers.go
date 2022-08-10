package jsonHelpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"unsafe"

	"github.com/mjarkk/jsonschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// RFC3339Nano is a time.Time that json (Un)Marshals from & to RFC3339 nano
type RFC3339Nano time.Time

// JSONSchemaDescribe implements jsonschema.Describe
func (RFC3339Nano) JSONSchemaDescribe() jsonschema.Property {
	minLen := uint(10)
	return jsonschema.Property{
		Title: "RFC3339 time string",
		Description: "This field is a RFC3339 (nano) time string, " +
			"RFC3339 is basicly an extension of ISO 8601 so that should also be fine here",
		Type: jsonschema.PropertyTypeString,
		Examples: []json.RawMessage{
			[]byte("\"2019-10-12T07:20:50.52Z\""),
			[]byte("\"2019-10-12T14:20:50.52+07:00\""),
		},
		MinLength: &minLen,
	}
}

// UnmarshalJSON transforms a RFC3339 string into *a
func (t *RFC3339Nano) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*t = RFC3339Nano(parsedTime)

	return nil
}

// MarshalJSON transforms a into t into a RFC3339 time string
func (t RFC3339Nano) MarshalJSON() ([]byte, error) {
	timeString := t.Time().Format(time.RFC3339Nano)
	return json.Marshal(timeString)
}

// UnmarshalBSONValue implements bson.ValueUnmarshaler
// by default RFC3339Nano is transformed to a empty map so here we fix that
func (t *RFC3339Nano) UnmarshalBSONValue(valueType bsontype.Type, data []byte) error {
	switch valueType {
	case bsontype.Null, bsontype.Undefined:
		// Do not set the value
		return nil
	case bsontype.DateTime:
		// Just continue
	default:
		return errors.New("expected bson datetime but got " + valueType.String())
	}

	timeInt, _, ok := bsoncore.ReadDateTime(data)
	if !ok {
		return errors.New("unable to parse bson datetime")
	}
	*t = RFC3339Nano(primitive.DateTime(timeInt).Time())
	return nil
}

// MarshalBSONValue implements bson.ValueMarshaler
// by default RFC3339Nano is transformed to a empty map so here we fix that
func (t RFC3339Nano) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(t.Time())
}

// Time returns the underlaying time object
func (t RFC3339Nano) Time() time.Time {
	return time.Time(t)
}

// ToPtr creates a pointer to t
// This is handy when you want to add a inline time to a struct field
func (t RFC3339Nano) ToPtr() *RFC3339Nano {
	return &t
}

// Format formats the time according to the format string
func (t *RFC3339Nano) Format(format string) string {
	if t == nil {
		return ""
	}
	return t.Time().Format(format)
}

// PhoneNumber represens a phone number
type PhoneNumber struct {
	IsLocal          bool   // 06 12345678
	HasCountryPrefix bool   // +31 6 12345678
	Number           uint64 // 612345678 (basically the number without any formatting)
}

// String converts the phone number into a string
func (n PhoneNumber) String() string {
	resp := ""
	if n.HasCountryPrefix {
		resp += "+"
	} else if n.IsLocal {
		resp += "0"
	}
	return resp + strconv.FormatUint(n.Number, 10)
}

// JSONSchemaDescribe implements jsonschema.Describe
func (PhoneNumber) JSONSchemaDescribe() jsonschema.Property {
	minLen := uint(3)
	return jsonschema.Property{
		Title:       "Phone number",
		Description: "This field can contain any phone number",
		Type:        jsonschema.PropertyTypeString,
		Examples: []json.RawMessage{
			[]byte("\"0612345678\""),
			[]byte("\"06 12345678\""),
			[]byte("\"+31 - 6 - 1234 - 5678\""),
		},
		MinLength: &minLen,
	}
}

// MarshalJSON transforms a phonenumber into a json string
func (n PhoneNumber) MarshalJSON() ([]byte, error) {
	resp := []byte{'"'}
	if n.HasCountryPrefix {
		resp = append(resp, '+')
	} else if n.IsLocal {
		resp = append(resp, '0')
	}
	return append(strconv.AppendUint(resp, n.Number, 10), '"'), nil
}

// ErrInvalidPhoneNumber is the error returned if the input phone is incorrect
var ErrInvalidPhoneNumber = errors.New("invalid phone number")

// UnmarshalJSON transforms reads a string and converts it into a
func (n *PhoneNumber) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return ErrInvalidPhoneNumber
	}

	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}

		hasCountryPrefix := false

		// Filter out the none numbers from the string
		b := []byte(s)
		for i := len(b) - 1; i >= 0; i-- {
			c := b[i]
			if c == '+' {
				hasCountryPrefix = true
			} else if c >= '0' && c <= '9' {
				// + should only appear before the first number if it apears after a + it's incorrect
				// as we loop backwards we can undo the changes made by c == '+'
				hasCountryPrefix = false
				continue
			}

			b = append(b[:i], b[i+1:]...)
		}

		if len(b) > 15 || len(b) < 3 {
			return ErrInvalidPhoneNumber
		}

		// We can parse s here as changes made to b are directly reflected on s
		s = *(*string)(unsafe.Pointer(&b))
		nr, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return ErrInvalidPhoneNumber
		}

		*n = PhoneNumber{
			HasCountryPrefix: hasCountryPrefix,
			IsLocal:          b[0] == '0',
			Number:           nr,
		}
	} else {
		var s uint64
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}

		if s < 100 {
			return ErrInvalidPhoneNumber
		}

		*n = PhoneNumber{
			IsLocal: true,
			Number:  s,
		}
	}

	return nil
}

// DriversLicense contains a European drivers license
//
// Why a 4 byte array?
//   - Fixed byte arrays are much faster than strings
//   - The drivers license code is never longer than 4 bytes
type DriversLicense [4]byte

// NewDriversLicense creates a DriversLicense from a string
func NewDriversLicense(name string) DriversLicense {
	resp := DriversLicense{0, 0, 0, 0}
	copy(resp[:], name)
	return resp
}

// DriversLicenses contains all the EU drivers licenses
// List obtained from https://en.wikipedia.org/wiki/European_driving_licence
var DriversLicenses = []DriversLicense{
	NewDriversLicense("AM"),
	NewDriversLicense("A1"),
	NewDriversLicense("A2"),
	NewDriversLicense("A"),
	NewDriversLicense("B"),
	NewDriversLicense("BE"),
	NewDriversLicense("B1"),
	NewDriversLicense("C1"),
	NewDriversLicense("C1E"),
	NewDriversLicense("C"),
	NewDriversLicense("CE"),
	NewDriversLicense("D1"),
	NewDriversLicense("D1E"),
	NewDriversLicense("D"),
	NewDriversLicense("DE"),
	NewDriversLicense("LK"),
	NewDriversLicense("S"),
	NewDriversLicense("BF17"),
	NewDriversLicense("L17"),
	NewDriversLicense("TR"),
	NewDriversLicense("TROL"),
	NewDriversLicense("H"),
	NewDriversLicense("TRAM"),
	NewDriversLicense("V"),
	NewDriversLicense("F"),
	NewDriversLicense("T"),
	NewDriversLicense("TM"),
	NewDriversLicense("K"),
	NewDriversLicense("L"),
	NewDriversLicense("G"),
	NewDriversLicense("W"),
}

// Strings converts a drivers license identifier into a string
func (dl DriversLicense) String() string {
	return string(bytes.TrimRightFunc(dl[:], func(r rune) bool {
		return r == 0
	}))
}

// JSONSchemaDescribe implements jsonschema.Describe
func (DriversLicense) JSONSchemaDescribe() jsonschema.Property {
	return jsonschema.Property{
		Title:       "EU Drivers license",
		Description: "Describes a EU drivers license, for more information see https://en.wikipedia.org/wiki/European_driving_licence",
		Type:        jsonschema.PropertyTypeString,
		Enum: func() []json.RawMessage {
			var enums []json.RawMessage
			for _, v := range DriversLicenses {
				enums = append(enums, json.RawMessage(`"`+v.String()+`"`))
			}
			return enums
		}(),
	}
}

// MarshalJSON transforms a drivers license into a json string
func (dl DriversLicense) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dl.String() + `"`), nil
}

// ErrInvalidDriversLicense is returned when a drivers license is invalid
var ErrInvalidDriversLicense = errors.New("invalid drivers license")

// UnmarshalJSON transforms reads a string and converts it into a drivers license
func (dl *DriversLicense) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || len(b) > 6 {
		// Drivers license needs to be at least 1 charcater and max 4
		// Note that b includes the quotes so we need to account for that in the length check
		return ErrInvalidDriversLicense
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("expected a string but got something else")
	}

	key := [4]byte{0, 0, 0, 0}
	copy(key[:], b[1:len(b)-1])

	// search for a matching drivers license
	for _, entry := range DriversLicenses {
		if entry == key {
			*dl = entry
			return nil
		}
	}

	// Maybe there where lowercase letters, lets correct them and search again for a matching drivers license
	for i := 0; i < 4; i++ {
		c := key[i]
		if c >= 'a' && c <= 'z' {
			// Convert to uppercase
			key[i] = c - 'a' - 'A'
		}
	}

	// search for a matching drivers license
	for _, entry := range DriversLicenses {
		if entry == key {
			*dl = entry
			return nil
		}
	}

	// :( we did not find a matching drivers license

	return ErrInvalidDriversLicense
}
