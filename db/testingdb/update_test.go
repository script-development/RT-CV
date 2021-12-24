package testingdb

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	testDB := NewDB()

	mockData := NewMockuser()
	Equal(t, "users", mockData.CollectionName())

	// Insert dummy data
	err := testDB.Insert(mockData)
	NoError(t, err)
	documentsCount, err := testDB.Count(mockData, nil)
	NoError(t, err)
	Equal(t, uint64(1), documentsCount)

	// Create a new mockuser so we don't update the value behind the pointer
	// what might cause a fake positive when checking the updated data in the database
	newMockData := NewMockuser()
	realname := "John Doe"
	newMockData.Realname = &realname
	newMockData.ID = mockData.ID

	err = testDB.UpdateByID(newMockData)
	NoError(t, err)

	// Check if the data in the database is actually replaced
	collectionData := testDB.getCollectionFromEntry(newMockData).data
	Equal(t, 1, len(collectionData))
	firstItem := collectionData[0].(*MockUser)
	NotNil(t, firstItem.Realname)
	Equal(t, realname, *firstItem.Realname)
}
