package dbHelpers

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMergeFilters(t *testing.T) {
	filters := MergeFilters()
	NotNil(t, filters)

	filters = MergeFilters(nil)
	NotNil(t, filters)

	filters = MergeFilters(bson.M{}, nil)
	NotNil(t, filters)

	filters = MergeFilters(nil, bson.M{})
	NotNil(t, filters)

	filters = MergeFilters(bson.M{"key": "value"}, bson.M{"other": "value"})
	Equal(t, bson.M{"key": "value", "other": "value"}, filters)

	filters = MergeFilters(bson.M{"key": "value"}, bson.M{"key": "value should be overwritten"})
	Equal(t, bson.M{"key": "value should be overwritten"}, filters)
}
