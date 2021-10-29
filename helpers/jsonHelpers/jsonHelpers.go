package jsonHelpers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mjarkk/jsonschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// RFC3339Nano is a time.Time that json (Un)Marshals from & to RFC3339 nano
type RFC3339Nano time.Time

// JSONSchemaDescribe implements schema.Describe
func (RFC3339Nano) JSONSchemaDescribe() jsonschema.Property {
	minLen := uint(10)
	return jsonschema.Property{
		Title: "RFC3339 time string",
		Description: "This field is a RFC3339 (nano) time string, " +
			"RFC3339 is basicly an extension of ISO 8601 so that should also be fine here",
		Type: jsonschema.PropertyTypeString,
		Examples: []interface{}{
			"2019-10-12T07:20:50.52Z",
			"2019-10-12T14:20:50.52+07:00",
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
