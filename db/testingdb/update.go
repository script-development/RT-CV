package testingdb

import (
	"github.com/script-development/RT-CV/db/dbInterfaces"
)

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
