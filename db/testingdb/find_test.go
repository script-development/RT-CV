package testingdb

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindOneWithoutFilters(t *testing.T) {
	testDB := NewDB()

	mockData := NewMockuser()
	Equal(t, "users", mockData.CollectionName())

	err := testDB.Insert(mockData)
	NoError(t, err)

	foundResult := MockUser{}
	err = testDB.FindOne(&foundResult, bson.M{})
	NoError(t, err)
	Equal(t, mockData.ID, foundResult.ID)
}

func TestFindWithoutFilters(t *testing.T) {
	testDB := NewDB()

	mockData := NewMockuser()
	Equal(t, "users", mockData.CollectionName())

	err := testDB.Insert(mockData)
	NoError(t, err)

	foundResults := []MockUser{}
	err = testDB.Find(&MockUser{}, &foundResults, bson.M{})
	NoError(t, err)
	Len(t, foundResults, 1)
	Equal(t, mockData.ID, foundResults[0].ID)

	foundResultsPtrs := []*MockUser{}
	err = testDB.Find(&MockUser{}, &foundResultsPtrs, bson.M{})
	NoError(t, err)
	Len(t, foundResultsPtrs, 1)
	Equal(t, mockData.ID, foundResultsPtrs[0].ID)
}
