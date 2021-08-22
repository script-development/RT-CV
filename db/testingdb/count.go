package testingdb

import "github.com/script-development/RT-CV/db"

// Count returns the number of documents in the collection of entity
func (c *TestConnection) Count(entity db.Entry) int {
	c.m.Lock()
	defer c.m.Unlock()

	collection, ok := c.collections[entity.CollectionName()]
	if !ok {
		return 0
	}

	return len(collection.data)
}
