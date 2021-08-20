package testingdb

import (
	"testing"

	"github.com/script-development/RT-CV/db/dbInterfaces"
	. "github.com/stretchr/testify/assert"
)

type MockUser struct {
	dbInterfaces.M `bson:",inline"`
	Realname       *string `bson:"real_name,omitempty"`
	Username       string
}

func (*MockUser) CollectionName() string {
	return "users"
}

func NewMockuser() *MockUser {
	return &MockUser{
		M:        dbInterfaces.NewM(),
		Realname: nil,
		Username: "Piet",
	}
}

type MockPost struct {
	dbInterfaces.M
	Title   string
	Content string
}

func (*MockPost) CollectionName() string {
	return "posts"
}

func TestNewDB(t *testing.T) {
	testDB := NewDB()
	NotNil(t, testDB)
}
