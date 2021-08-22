package testingdb

import (
	"testing"

	"github.com/script-development/RT-CV/db"
	. "github.com/stretchr/testify/assert"
)

type MockUser struct {
	db.M     `bson:",inline"`
	Realname *string `bson:"real_name,omitempty"`
	Username string
}

func (*MockUser) CollectionName() string {
	return "users"
}

func NewMockuser() *MockUser {
	return &MockUser{
		M:        db.NewM(),
		Realname: nil,
		Username: "Piet",
	}
}

type MockPost struct {
	db.M
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
