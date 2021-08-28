package jsonHelpers

import (
	"encoding/json"
	"time"

	"github.com/script-development/RT-CV/helpers/schema"
)

// RFC3339Nano is a time.Time that json (Un)Marshals from & to RFC3339 nano
type RFC3339Nano time.Time

// JSONSchemaDescribe implments schema.Describe
func (RFC3339Nano) JSONSchemaDescribe() schema.Property {
	minLen := uint(10)
	return schema.Property{
		Title:       "RFC3339 time string",
		Description: "This field is a RFC3339 (nano) time string that requires an integer to work",
		Type:        schema.PropertyTypeString,
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

// Time returns the underlaying time object
func (t RFC3339Nano) Time() time.Time {
	return time.Time(t)
}

// ToPtr creates a pointer to t
// This is handy when you want to add a inline time to a struct field
func (t RFC3339Nano) ToPtr() *RFC3339Nano {
	return &t
}
