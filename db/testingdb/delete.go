package testingdb

import (
	"github.com/script-development/RT-CV/db/dbInterfaces"
)

func (c *TestConnection) DeleteByID(e dbInterfaces.Entry) error {
	c.m.Lock()
	defer c.m.Unlock()

	eId := e.GetID()

	collection := c.getCollectionFromEntry(e)
	for i, entry := range collection.data {
		if entry.GetID() == eId {
			collection.data = append(collection.data[:i], collection.data[i+1:]...)
			c.setCollection(collection)
			break
		}
	}

	return nil
}
