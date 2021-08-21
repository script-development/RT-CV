package testingdb

import (
	"sync"

	"github.com/script-development/RT-CV/db/dbInterfaces"
)

// TestConnection is the struct that implements dbInterfaces.Connection
type TestConnection struct {
	m           sync.Mutex
	collections map[string]Collection
}

// NewDB returns a testing database connection that is compatible with dbInterfaces.Connection
func NewDB() *TestConnection {
	return &TestConnection{
		collections: map[string]Collection{},
	}
}

// Collection contains all the data for a collection
type Collection struct {
	name string
	data []dbInterfaces.Entry
}

// RegisterEntries implements dbInterfaces.Connection
func (c *TestConnection) RegisterEntries(entries ...dbInterfaces.Entry) {
	// We don't need to implement this function
}

func (c *TestConnection) getCollectionFromEntry(e dbInterfaces.Entry) Collection {
	collectionName := e.CollectionName()
	v, ok := c.collections[collectionName]
	if ok {
		return v
	}

	return Collection{
		name: collectionName,
		data: []dbInterfaces.Entry{},
	}
}

func (c *TestConnection) setCollection(collection Collection) {
	c.collections[collection.name] = collection
}
