package testingdb

import (
	"errors"
	"reflect"

	"github.com/script-development/RT-CV/db/dbInterfaces"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *TestConnection) FindOne(placeInto dbInterfaces.Entry, filters bson.M) error {
	itemsFilter := newFilter(placeInto.DefaultFindFilters(), filters)

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

	return errors.New("no document found")
}

func (c *TestConnection) Find(base dbInterfaces.Entry, results interface{}, filters bson.M) error {
	itemsFilter := newFilter(base.DefaultFindFilters(), filters)

	c.m.Lock()
	defer c.m.Unlock()

	resultRefl := reflect.ValueOf(results).Elem()
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
