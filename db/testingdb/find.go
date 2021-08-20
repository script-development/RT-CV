package testingdb

import (
	"errors"
	"reflect"

	"github.com/script-development/RT-CV/db/dbHelpers"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *TestConnection) FindOne(placeInto dbInterfaces.Entry, filters bson.M) error {
	mergedFilters := dbHelpers.MergeFilters(placeInto.DefaultFindFilters(), filters)
	if len(mergedFilters) > 0 {
		panic("TODO impl filters")
	}

	c.m.Lock()
	defer c.m.Unlock()

	for _, item := range c.getCollectionFromEntry(placeInto).data {

		// We use elem here to get passed the pointer into the underlaying data
		placeIntoRefl := reflect.ValueOf(placeInto).Elem()
		placeIntoRefl.Set(reflect.ValueOf(item).Elem())
		return nil
	}

	return errors.New("no document found")
}

func (c *TestConnection) Find(base dbInterfaces.Entry, results interface{}, filters bson.M) error {
	mergedFilters := dbHelpers.MergeFilters(base.DefaultFindFilters(), filters)
	if len(mergedFilters) > 0 {
		panic("TODO impl filters")
	}

	c.m.Lock()
	defer c.m.Unlock()

	resultRefl := reflect.ValueOf(results).Elem()
	if resultRefl.Kind() != reflect.Slice {
		return errors.New("requires pointer to slice as results argument")
	}

	resultsSliceContentType := resultRefl.Type().Elem()
	placePtr := resultsSliceContentType.Kind() == reflect.Ptr

	for _, item := range c.getCollectionFromEntry(base).data {
		itemRefl := reflect.ValueOf(item)
		if placePtr {
			resultRefl = reflect.Append(resultRefl, itemRefl)
		} else {
			resultRefl = reflect.Append(resultRefl, itemRefl.Elem())
		}
	}

	reflect.ValueOf(results).Elem().Set(resultRefl)

	return nil
}
