package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	// Returns err == mongo.ErrNoDocuments if no documents where found
	FindOne(result Entry, filters bson.M, opts ...FindOptions) error

	// Find finds multiple entries in the database
	// The entry argument is to determain on which collection we execute the query
	Find(entry Entry, results any, filters bson.M, opts ...FindOptions) error

	// Insert inserts an entry into the database
	Insert(data ...Entry) error

	// UpdateID updates an entry in the database
	UpdateByID(data Entry) error

	// DeleteByID deletes an entry from the database
	DeleteByID(entry Entry, ids ...primitive.ObjectID) error

	// Count counts the number of documents in the database for the specific filter
	// If filter is nil the number of all the documents is returned
	Count(entry Entry, filter bson.M) (uint64, error)
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
	// Get the _id field of the entry
	GetID() primitive.ObjectID

	// Set the _id field of the entry
	SetID(primitive.ObjectID)

	// CollectionName should yield the collection name for the entre
	CollectionName() string

	// DefaultFindFilters can return a default filter used in find queries
	// If nil is returned this is not used
	DefaultFindFilters() bson.M

	// Indexes returns the indexes for the entry
	// If nil is returned no more indexes will be set
	// Note that by default the there is always an index of the _id field
	Indexes() []mongo.IndexModel
}

// M is a struct that adds an _id field and implements from dbInterfaces.Entry:
// - GetID
// - SetID
// - DefaultFindFilters
type M struct {
	ID primitive.ObjectID `bson:"_id" json:"id" description:"The unique id of the entry in the MongoDB ObjectId format, for more info see: https://docs.mongodb.com/manual/reference/method/ObjectId/"`
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

// Indexes implements Entry
func (*M) Indexes() []mongo.IndexModel {
	return nil
}
