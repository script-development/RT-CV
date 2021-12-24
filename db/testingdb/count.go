package testingdb

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Count returns the number of documents in the collection of entity
func (c *TestConnection) Count(entry db.Entry, filter bson.M) (uint64, error) {
	c.m.Lock()
	defer c.m.Unlock()

	collection := c.getCollectionFromEntry(entry)
	if len(filter) == 0 {
		// Take the easy route
		return uint64(len(collection.data)), nil
	}

	itemsFilter := newFilter(filter)
	var count uint64
	for _, item := range collection.data {
		if itemsFilter.matches(item) {
			count++
		}
	}

	return count, nil
}
