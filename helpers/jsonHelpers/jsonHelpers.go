package jsonHelpers

import (
	"encoding/json"
	"time"
)

// RFC3339Nano is a time.Time that json (Un)Marshals from & to RFC3339 nano
type RFC3339Nano time.Time

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
