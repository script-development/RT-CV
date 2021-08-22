package testingdb

import (
	"sync"

	"github.com/script-development/RT-CV/db"
)

// TestConnection is the struct that implements db.Connection
type TestConnection struct {
	m           sync.Mutex
	collections map[string]Collection
}

// NewDB returns a testing database connection that is compatible with db.Connection
func NewDB() *TestConnection {
	return &TestConnection{
		collections: map[string]Collection{},
	}
}

// Collection contains all the data for a collection
type Collection struct {
	name string
	data []db.Entry
}

// RegisterEntries implements db.Connection
func (*TestConnection) RegisterEntries(...db.Entry) {
	// We don't need to implement this function
}

func (c *TestConnection) getCollectionFromEntry(e db.Entry) Collection {
	collectionName := e.CollectionName()
	v, ok := c.collections[collectionName]
	if ok {
		return v
	}

	return Collection{
		name: collectionName,
		data: []db.Entry{},
	}
}

func (c *TestConnection) setCollection(collection Collection) {
	c.collections[collection.name] = collection
}
