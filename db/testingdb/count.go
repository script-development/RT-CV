package testingdb

import "github.com/script-development/RT-CV/db/dbInterfaces"

func (c *TestConnection) Count(entity dbInterfaces.Entry) int {
	c.m.Lock()
	defer c.m.Unlock()

	collection, ok := c.collections[entity.CollectionName()]
	if !ok {
		return 0
	}

	return len(collection.data)
}
