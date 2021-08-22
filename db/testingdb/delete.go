package testingdb

import (
	"github.com/script-development/RT-CV/db"
)

// DeleteByID deletes a document by it's ID
func (c *TestConnection) DeleteByID(e db.Entry) error {
	c.m.Lock()
	defer c.m.Unlock()

	eID := e.GetID()

	collection := c.getCollectionFromEntry(e)
	for i, entry := range collection.data {
		if entry.GetID() == eID {
			collection.data = append(collection.data[:i], collection.data[i+1:]...)
			c.setCollection(collection)
			break
		}
	}

	return nil
}
