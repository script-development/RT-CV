package testingdb

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	testDB := NewDB()

	mockData := NewMockuser()
	Equal(t, "users", mockData.CollectionName())

	// Insert dummy data
	err := testDB.Insert(mockData)
	NoError(t, err)
	documentsCount, _ := testDB.Count(mockData, nil)
	Equal(t, uint64(1), documentsCount)

	// Delete entry and check if the collection is now empty
	err = testDB.DeleteByID(&MockUser{}, mockData.ID)
	NoError(t, err)
	documentsCount, _ = testDB.Count(mockData, nil)
	Equal(t, uint64(0), documentsCount)

	// Should result in no panics/errors if there is nothing to delete
	err = testDB.DeleteByID(&MockUser{}, mockData.ID)
	NoError(t, err)
	documentsCount, _ = testDB.Count(mockData, nil)
	Equal(t, uint64(0), documentsCount)
}
