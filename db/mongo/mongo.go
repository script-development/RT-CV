package mongo

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/dbHelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectToDB connects to a mongodb database based on a shell variable ($MONGODB_URI)
func ConnectToDB() db.Connection {
	log.Info("Connecting to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatalf("Mongo new client failed, error: " + err.Error())
	}

	ctx, ctxCancel := context.WithTimeout(dbHelpers.Ctx(), 10*time.Second)
	err = client.Connect(ctx)
	ctxCancel()
	if err != nil {
		log.Fatal("Mongo client connect failed, error: " + err.Error())
	}

	ctx, ctxCancel = context.WithTimeout(dbHelpers.Ctx(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	ctxCancel()
	if err != nil {
		log.Fatal("Mongo client ping failed, error: " + err.Error())
	}

	mongoConnection := client.Database(os.Getenv("MONGODB_DATABASE"))
	log.Info("Connected to database")
	var conn db.Connection = &Connection{
		db: mongoConnection,
	}
	return conn
}

// Connection is a connection to a mongo database
type Connection struct {
	db *mongo.Database
}

// FindOne finds a single entry based on the filter
func (c *Connection) FindOne(e db.Entry, filter bson.M, optionalOpts ...db.FindOptions) error {
	opts := db.FindOptions{}
	if len(optionalOpts) > 0 {
		opts = optionalOpts[0]
	}

	queryFilters := filter
	if !opts.NoDefaultFilters {
		dbHelpers.MergeFilters(e.DefaultFindFilters(), filter)
	}

	res := c.collection(e).FindOne(dbHelpers.Ctx(), queryFilters)
	err := res.Err()
	if err != nil {
		return err
	}
	err = res.Decode(e)
	return err
}

// Find finds entries based on the filter
func (c *Connection) Find(e db.Entry, results interface{}, filter bson.M, optionalOpts ...db.FindOptions) error {
	opts := db.FindOptions{}
	if len(optionalOpts) > 0 {
		opts = optionalOpts[0]
	}

	queryFilters := filter
	if !opts.NoDefaultFilters {
		dbHelpers.MergeFilters(e.DefaultFindFilters(), filter)
	}

	cur, err := c.collection(e).Find(dbHelpers.Ctx(), queryFilters)
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
func (c *Connection) Insert(e ...db.Entry) error {
	for idx, entry := range e {
		if entry.GetID().IsZero() {
			e[idx].SetID(primitive.NewObjectID())
		}
	}

	switch len(e) {
	case 0:
		return nil
	case 1:
		_, err := c.collection(e[0]).InsertOne(dbHelpers.Ctx(), e[0])
		return err
	default:
		// Convert e to a slice of interface{}
		// Fixes: panic: interface conversion: interface {} is []db.Entry, not []interface {}
		eAsInterf := []interface{}{}
		for _, entry := range e {
			eAsInterf = append(eAsInterf, entry)
		}

		_, err := c.collection(e[0]).InsertMany(dbHelpers.Ctx(), eAsInterf)
		return err
	}
}

// UpdateByID updates an entry by its id
func (c *Connection) UpdateByID(e db.Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).ReplaceOne(dbHelpers.Ctx(), bson.M{"_id": id}, e)
	return err
}

// DeleteByID deletes an entry by its id
func (c *Connection) DeleteByID(e db.Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).DeleteOne(dbHelpers.Ctx(), bson.M{"_id": id})
	return err
}

// Count counts documents in a collection
func (c *Connection) Count(entry db.Entry, filter bson.M) (uint64, error) {
	if filter == nil {
		filter = bson.M{}
	}
	count, err := c.collection(entry).CountDocuments(dbHelpers.Ctx(), filter)
	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}

func (c *Connection) collection(entry db.Entry) *mongo.Collection {
	return c.db.Collection(entry.CollectionName())
}

// RegisterEntries creates a collection for every entry
func (c *Connection) RegisterEntries(entries ...db.Entry) {
	log.Info("Checking if all db collections exist")

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
			log.Infof("Creating collection %s", collectionName)
			err = c.db.CreateCollection(dbHelpers.Ctx(), collectionName)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		indexes := entry.Indexes()
		if len(indexes) > 0 {
			_, err = c.db.Collection(collectionName).Indexes().CreateMany(context.Background(), indexes)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	log.Info("Database collection check succeeded")
}
