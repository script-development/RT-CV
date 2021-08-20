package dbHelpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// MergeFilters merges a list of filters into a single filter
func MergeFilters(filtersList ...bson.M) bson.M {
	res := filtersList[0]

	for _, filters := range filtersList[1:] {
		if res == nil {
			res = filters
			continue
		}

		for key, value := range filters {
			res[key] = value
		}
	}

	if res == nil {
		return bson.M{}
	}
	return res
}

// Ctx
func Ctx() context.Context {
	return context.Background()
}
