package dbInterfaces

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FakeConn struct {
	Defer Connection
}

// Connection is a abstract interface for a database connection
// There are 2 main implementations of this:
// - MongoConnection (For the MongoDB driver)
// - TestConnection (For a fake temp database)
type Connection interface {
	RegisterEntries(entries ...Entry)
	FindOne(result Entry, filters bson.M) error
	Find(entry Entry, results interface{}, filters bson.M) error
	Insert(data Entry) error
	UpdateByID(data Entry) error
	DeleteByID(data Entry) error
}

// Entry are the functions required to put/get things in/from the database
// To implement use:
//
// type User struct {
//     // Adds the _id field and implements the remaining functions from Entry
//     M `bson:",inline"`
//
//     Username string
// }
// func (*User) CollectionName() {
//     return "users"
// }
type Entry interface {
	GetID() primitive.ObjectID
	SetID(primitive.ObjectID)
	CollectionName() string
	DefaultFindFilters() bson.M
}

// M is a struct that adds an _id field and implements from dbInterfaces.Entry:
// - GetID
// - SetID
// - DefaultFindFilters
type M struct {
	ID primitive.ObjectID `bson:"_id"`
}

func NewM() M {
	return M{
		ID: primitive.NewObjectID(),
	}
}

// GetID implements Entry
func (m *M) GetID() primitive.ObjectID {
	return m.ID
}

// SetID implements Entry
func (m *M) SetID(id primitive.ObjectID) {
	m.ID = id
}

// DefaultFindFilters implements Entry
func (m *M) DefaultFindFilters() bson.M {
	return nil
}
