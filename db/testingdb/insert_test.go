package testingdb

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	testDB := NewDB()

	mockData := NewMockuser()
	Equal(t, "users", mockData.CollectionName())

	err := testDB.Insert(mockData)

	v, ok := testDB.collections["users"]
	True(t, ok)
	NotNil(t, v)
	Equal(t, "users", v.name)
	Len(t, v.data, 1)
	Equal(t, mockData.GetID(), v.data[0].GetID())

	NoError(t, err)
}
