package testingdb

import (
	"github.com/script-development/RT-CV/db/dbInterfaces"
)

// UpdateByID updates a document in the database by its ID
func (c *TestConnection) UpdateByID(updateData dbInterfaces.Entry) error {
	c.m.Lock()
	defer c.m.Unlock()

	updateDataID := updateData.GetID()
	collection := c.getCollectionFromEntry(updateData)

	for i, entry := range collection.data {
		if entry.GetID() == updateDataID {
			collection.data[i] = updateData
			c.setCollection(collection)
			break
		}
	}

	return nil
}
