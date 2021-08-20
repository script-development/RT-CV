package testingdb

import (
	"sync"

	"github.com/script-development/RT-CV/db/dbInterfaces"
)

func NewDB() *TestConnection {
	return &TestConnection{
		collections: map[string]Collection{},
	}
}

type TestConnection struct {
	m           sync.Mutex
	collections map[string]Collection
}

type Collection struct {
	name string
	data []dbInterfaces.Entry
}

func (c *TestConnection) RegisterEntries(entries ...dbInterfaces.Entry) {
	// This function doesn't have to be implemented
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
