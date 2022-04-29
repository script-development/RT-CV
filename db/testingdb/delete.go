package testingdb

import (
	"github.com/script-development/RT-CV/db"
)

// DeleteByID deletes a document by it's ID
func (c *TestConnection) DeleteByID(e ...db.Entry) error {
	c.m.Lock()
	defer c.m.Unlock()

	if len(e) == 0 {
		return nil
	}

	collection := c.getCollectionFromEntry(e[0])
	for _, entry := range e {
		eID := entry.GetID()
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
