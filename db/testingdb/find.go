package testingdb

import (
	"errors"
	"reflect"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/dbHelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOne finds one document in the collection of placeInto
// The result can be filtered using filters
// The filters should work equal to MongoDB filters (https://docs.mongodb.com/manual/tutorial/query-documents/)
// tough this might miss features compared to mongoDB's filters
func (c *TestConnection) FindOne(placeInto db.Entry, filters bson.M, optionalOpts ...db.FindOptions) error {
	opts := db.FindOptions{}
	if len(optionalOpts) > 0 {
		opts = optionalOpts[0]
	}

	queryFilters := filters
	if !opts.NoDefaultFilters {
		dbHelpers.MergeFilters(placeInto.DefaultFindFilters(), filters)
	}
	itemsFilter := newFilter(queryFilters)

	c.m.Lock()
	defer c.m.Unlock()

	for _, item := range c.getCollectionFromEntry(placeInto).data {
		if !itemsFilter.matches(item) {
			continue
		}

		// We use elem here to get passed the pointer into the underlaying data
		placeIntoRefl := reflect.ValueOf(placeInto).Elem()
		placeIntoRefl.Set(reflect.ValueOf(item).Elem())
		return nil
	}

	return mongo.ErrNoDocuments
}

// Find finds documents in the collection of the base
// The results can be filtered using filters
// The filters should work equal to MongoDB filters (https://docs.mongodb.com/manual/tutorial/query-documents/)
// tough this might miss features compared to mongoDB's filters
func (c *TestConnection) Find(base db.Entry, results interface{}, filters bson.M, optionalOpts ...db.FindOptions) error {
	opts := db.FindOptions{}
	if len(optionalOpts) > 0 {
		opts = optionalOpts[0]
	}

	queryFilters := filters
	if !opts.NoDefaultFilters {
		dbHelpers.MergeFilters(base.DefaultFindFilters(), filters)
	}
	itemsFilter := newFilter(queryFilters)

	c.m.Lock()
	defer c.m.Unlock()

	resultRefl := reflect.ValueOf(results)
	if resultRefl.Kind() != reflect.Ptr {
		return errors.New("requires pointer to slice as results argument")
	}

	resultRefl = resultRefl.Elem()
	if resultRefl.Kind() != reflect.Slice {
		return errors.New("requires pointer to slice as results argument")
	}

	resultsSliceContentType := resultRefl.Type().Elem()
	resultIsSliceOfPtrs := resultsSliceContentType.Kind() == reflect.Ptr

	for _, item := range c.getCollectionFromEntry(base).data {
		if !itemsFilter.matches(item) {
			continue
		}

		itemRefl := reflect.ValueOf(item)
		if resultIsSliceOfPtrs {
			resultRefl = reflect.Append(resultRefl, itemRefl)
		} else {
			resultRefl = reflect.Append(resultRefl, itemRefl.Elem())
		}
	}

	reflect.ValueOf(results).Elem().Set(resultRefl)

	return nil
}
