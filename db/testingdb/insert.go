package testingdb

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert inserts an item into the database
// Implements db.Connection
func (c *TestConnection) Insert(data ...db.Entry) error {
	c.m.Lock()
	defer c.m.Unlock()

	return c.UnsafeInsert(data...)
}

// UnsafeInsert inserts data directly into the database
// !!Be aware!!:
// - Not thread safe
// - Assumes the all entries are of same type / collection
func (c *TestConnection) UnsafeInsert(entries ...db.Entry) error {
	if len(entries) == 0 {
		return nil
	}
	for _, entry := range entries {
		if entry.GetID().IsZero() {
			entry.SetID(primitive.NewObjectID())
		}
	}

	collection := c.getCollectionFromEntry(entries[0])
	collection.data = append(collection.data, entries...)
	c.setCollection(collection)
	return nil
}
