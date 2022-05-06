package testingdb

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteByID deletes a document by it's ID
func (c *TestConnection) DeleteByID(entry db.Entry, ids ...primitive.ObjectID) error {
	c.m.Lock()
	defer c.m.Unlock()

	if len(ids) == 0 {
		return nil
	}

	collection := c.getCollectionFromEntry(entry)
	for _, eID := range ids {
		for i, collectionEntry := range collection.data {
			if collectionEntry.GetID() == eID {
				collection.data = append(collection.data[:i], collection.data[i+1:]...)
				c.setCollection(collection)
				break
			}
		}
	}

	return nil
}
