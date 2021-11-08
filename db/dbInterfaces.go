package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindOptions are options for the find operator in Connection
type FindOptions struct {
	// NoDefaultFilters does not include the default filters for the entry provided
	NoDefaultFilters bool
}

// Connection is a abstract interface for a database connection
// There are 2 main implementations of this:
// - MongoConnection (For the MongoDB driver)
// - TestConnection (For a fake temp database)
type Connection interface {
	// RegisterEntries tells the database to register the given entries
	// In the case of the mongodb database this means we'll create a collection for each entry
	RegisterEntries(entries ...Entry)

	// FindOne finds one entry inside the database
	FindOne(result Entry, filters bson.M, opts ...FindOptions) error

	// Find finds multiple entries in the database
	Find(entry Entry, results interface{}, filters bson.M, opts ...FindOptions) error

	// Insert inserts an entry into the database
	Insert(data ...Entry) error

	// UpdateID updates an entry in the database
	UpdateByID(data Entry) error

	// DeleteByID deletes an entry from the database
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
	ID primitive.ObjectID `bson:"_id" json:"id"`
}

// NewM returns a new instance of M
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
func (*M) DefaultFindFilters() bson.M {
	return nil
}
