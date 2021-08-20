package testingdb

import (
	"github.com/script-development/RT-CV/db/dbInterfaces"
)

func (c *TestConnection) Insert(data dbInterfaces.Entry) error {
	c.m.Lock()
	defer c.m.Unlock()

	collection := c.getCollectionFromEntry(data)
	collection.data = append(collection.data, data)
	c.setCollection(collection)

	return nil
}
