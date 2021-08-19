package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToDB() Connection {
	fmt.Println("Connecting to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel := context.WithTimeout(dbCtx(), 10*time.Second)
	err = client.Connect(ctx)
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel = context.WithTimeout(dbCtx(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	fmt.Println("Connected to database")
	var conn Connection = &MongoConnection{
		db: db,
	}
	return conn
}

type MongoConnection struct {
	db *mongo.Database
}

func (c *MongoConnection) FindOne(e Entry, filter bson.M) error {
	res := c.collection(e).FindOne(dbCtx(), mergeFilters(e.DefaultFindFilters(), filter))
	err := res.Err()
	if err != nil {
		return err
	}
	err = res.Decode(e)
	return err
}
func (c *MongoConnection) Find(e Entry, results interface{}, filter bson.M) error {
	cur, err := c.collection(e).Find(dbCtx(), mergeFilters(e.DefaultFindFilters(), filter))
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	err = cur.All(dbCtx(), results)
	return err
}
func (c *MongoConnection) Insert(e Entry) error {
	if e.GetID().IsZero() {
		e.SetID(primitive.NewObjectID())
	}

	_, err := c.collection(e).InsertOne(dbCtx(), e)
	return err
}

func (c *MongoConnection) UpdateByID(e Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).UpdateOne(dbCtx(), bson.M{"_id": id}, e)
	return err
}

func (c *MongoConnection) DeleteByID(e Entry) error {
	id := e.GetID()
	if id.IsZero() {
		return errors.New("cannot update item without id")
	}

	_, err := c.collection(e).DeleteOne(dbCtx(), bson.M{"_id": id})
	return err
}

func (c *MongoConnection) collection(entry Entry) *mongo.Collection {
	return c.db.Collection(entry.CollectionName())
}

func (c *MongoConnection) RegisterEntries(entries ...Entry) {
	fmt.Println("Checking if all db collections exist")

	names, err := c.db.ListCollectionNames(dbCtx(), bson.D{})
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
			c.db.CreateCollection(dbCtx(), collectionName)
		}
	}

	fmt.Println("Database collection check succeeded")
}
