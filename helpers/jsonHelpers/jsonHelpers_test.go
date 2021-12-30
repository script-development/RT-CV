package jsonHelpers

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRFC3339Nano(t *testing.T) {
	baseTime := time.Now()
	expectedJSONOutput := baseTime.Format(time.RFC3339Nano)

	marshaledTime, err := json.Marshal(RFC3339Nano(baseTime))
	NoError(t, err)
	Equal(t, `"`+expectedJSONOutput+`"`, string(marshaledTime))

	var parsedTime RFC3339Nano
	err = json.Unmarshal(marshaledTime, &parsedTime)
	NoError(t, err)

	parsedTimeAsTime := parsedTime.Time()
	False(t, parsedTimeAsTime.IsZero())

	Equal(t, baseTime.Year(), parsedTimeAsTime.Year())
	Equal(t, baseTime.Month(), parsedTimeAsTime.Month())
	Equal(t, baseTime.Day(), parsedTimeAsTime.Day())

	Equal(t, baseTime.Hour(), parsedTimeAsTime.Hour())
	Equal(t, baseTime.Minute(), parsedTimeAsTime.Minute())
	Equal(t, baseTime.Second(), parsedTimeAsTime.Second())

	type TestStruct struct {
		foo RFC3339Nano
	}
	bytes, err := bson.Marshal(TestStruct{RFC3339Nano(time.Now())})
	NoError(t, err)
	parsedTestStruct := TestStruct{}
	err = bson.Unmarshal(bytes, &parsedTestStruct)
	NoError(t, err)
}

func TestPhoneNumber(t *testing.T) {
	tests := []struct {
		input       string
		expects     PhoneNumber
		expectsJSON string
	}{
		{"+49123456789", PhoneNumber{HasCountryPrefix: true, Number: 49123456789}, "+49123456789"},
		{"+49 123 456 789", PhoneNumber{HasCountryPrefix: true, Number: 49123456789}, "+49123456789"},
		{"+49 - 123 - 456 - 789", PhoneNumber{HasCountryPrefix: true, Number: 49123456789}, "+49123456789"},
		{"0612345678", PhoneNumber{IsLocal: true, Number: 612345678}, "0612345678"},
		{"06 1234 5678", PhoneNumber{IsLocal: true, Number: 612345678}, "0612345678"},
		{"06 + 1234 + 5678", PhoneNumber{IsLocal: true, Number: 612345678}, "0612345678"},
	}

	for _, testCase := range tests {
		var phoneNumber PhoneNumber
		err := json.Unmarshal([]byte(`"`+testCase.input+`"`), &phoneNumber)
		NoError(t, err)
		Equal(t, testCase.expects, phoneNumber)

		marshaled, err := json.Marshal(phoneNumber)
		NoError(t, err)
		Equal(t, `"`+testCase.expectsJSON+`"`, string(marshaled))
	}
}
