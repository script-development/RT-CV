package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db/dbHelpers"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectToDB connects to a mongodb database based on a shell variable ($MONGODB_URI)
func ConnectToDB() dbInterfaces.Connection {
	fmt.Println("Connecting to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel := context.WithTimeout(dbHelpers.Ctx(), 10*time.Second)
	err = client.Connect(ctx)
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel = context.WithTimeout(dbHelpers.Ctx(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	fmt.Println("Connected to database")
	var conn dbInterfaces.Connection = &Connection{
		db: db,
	}
	return conn
}

// Connection is a connection to a mongo database
type Connection struct {
	db *mongo.Database
}

// FindOne finds a single entry based on the filter
func (c *Connection) FindOne(e dbInterfaces.Entry, filter bson.M) error {
	res := c.collection(e).FindOne(dbHelpers.Ctx(), dbHelpers.MergeFilters(e.DefaultFindFilters(), filter))
	err := res.Err()
	if err != nil {
		return err
	}
	err = res.Decode(e)
	return err
}

// Find finds entries based on the filter
func (c *Connection) Find(e dbInterfaces.Entry, results interface{}, filter bson.M) error {
	cur, err := c.collection(e).Find(dbHelpers.Ctx(), dbHelpers.MergeFilters(e.DefaultFindFilters(), filter))
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	err = cur.All(dbHelpers.Ctx(), results)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}

// Insert inserts an entry into the database
func (c *Connection) Insert(e dbInterfaces.Entry) error {
	if e.GetID().IsZero() {
		e.SetID(primitive.NewObjectID())
	}

	_, err := c.collection(e).InsertOne(dbHelpers.Ctx(), e)
	return err
}

// UpdateByID updates an entry by its id
func (c *Connection) UpdateByID(e dbInterfaces.Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).UpdateOne(dbHelpers.Ctx(), bson.M{"_id": id}, e)
	return err
}

// DeleteByID deletes an entry by its id
func (c *Connection) DeleteByID(e dbInterfaces.Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).DeleteOne(dbHelpers.Ctx(), bson.M{"_id": id})
	return err
}

func (c *Connection) collection(entry dbInterfaces.Entry) *mongo.Collection {
	return c.db.Collection(entry.CollectionName())
}

// RegisterEntries creates a collection for every entry
func (c *Connection) RegisterEntries(entries ...dbInterfaces.Entry) {
	fmt.Println("Checking if all db collections exist")

	names, err := c.db.ListCollectionNames(dbHelpers.Ctx(), bson.D{})
	if err != nil {
		log.Fatal(err.Error())
	}
	namesMap := map[string]bool{}
	for _, name := range names {
		namesMap[name] = true
	}
	for _, entry := range entries {
		collectionName := entry.CollectionName()

		if !namesMap[collectionName] {
			fmt.Println("Creating collection", collectionName)
			c.db.CreateCollection(dbHelpers.Ctx(), collectionName)
		}
	}

	fmt.Println("Database collection check succeeded")
}
