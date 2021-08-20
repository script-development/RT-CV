package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/db/dbhelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToDB() dbInterfaces.Connection {
	fmt.Println("Connecting to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel := context.WithTimeout(dbhelpers.Ctx(), 10*time.Second)
	err = client.Connect(ctx)
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel = context.WithTimeout(dbhelpers.Ctx(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	fmt.Println("Connected to database")
	var conn dbInterfaces.Connection = &MongoConnection{
		db: db,
	}
	return conn
}

type MongoConnection struct {
	db *mongo.Database
}

func (c *MongoConnection) FindOne(e dbInterfaces.Entry, filter bson.M) error {
	res := c.collection(e).FindOne(dbhelpers.Ctx(), dbhelpers.MergeFilters(e.DefaultFindFilters(), filter))
	err := res.Err()
	if err != nil {
		return err
	}
	err = res.Decode(e)
	return err
}
func (c *MongoConnection) Find(e dbInterfaces.Entry, results interface{}, filter bson.M) error {
	cur, err := c.collection(e).Find(dbhelpers.Ctx(), dbhelpers.MergeFilters(e.DefaultFindFilters(), filter))
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	err = cur.All(dbhelpers.Ctx(), results)
	return err
}
func (c *MongoConnection) Insert(e dbInterfaces.Entry) error {
	if e.GetID().IsZero() {
		e.SetID(primitive.NewObjectID())
	}

	_, err := c.collection(e).InsertOne(dbhelpers.Ctx(), e)
	return err
}

func (c *MongoConnection) UpdateByID(e dbInterfaces.Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).UpdateOne(dbhelpers.Ctx(), bson.M{"_id": id}, e)
	return err
}

func (c *MongoConnection) DeleteByID(e dbInterfaces.Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).DeleteOne(dbhelpers.Ctx(), bson.M{"_id": id})
	return err
}

func (c *MongoConnection) collection(entry dbInterfaces.Entry) *mongo.Collection {
	return c.db.Collection(entry.CollectionName())
}

func (c *MongoConnection) RegisterEntries(entries ...dbInterfaces.Entry) {
	fmt.Println("Checking if all db collections exist")

	names, err := c.db.ListCollectionNames(dbhelpers.Ctx(), bson.D{})
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
			c.db.CreateCollection(dbhelpers.Ctx(), collectionName)
		}
	}

	fmt.Println("Database collection check succeeded")
}
