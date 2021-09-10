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
