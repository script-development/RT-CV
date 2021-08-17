package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

func Ctx() context.Context {
	return context.Background()
}

func ConnectToDB() {
	fmt.Println("Connecting to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, ctxCancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	ctxCancel()
	if err != nil {
		log.Fatal(err.Error())
	}

	DB = client.Database(os.Getenv("MONGODB_DATABASE"))
	fmt.Println("Connected to database")
}

type Collection string

const (
	ApiKeys  = Collection("apiKeys")
	Profiles = Collection("profiles")
)

var AllCollections = []Collection{
	ApiKeys,
	Profiles,
}

func (c Collection) Collection() *mongo.Collection {
	return DB.Collection(string(c))
}

func InitDB() {
	fmt.Println("Checking if all db collections exist")

	names, err := DB.ListCollectionNames(Ctx(), bson.D{})
	if err != nil {
		log.Fatal(err.Error())
	}
	namesMap := map[string]bool{}
	for _, name := range names {
		namesMap[name] = true
	}
	for _, collection := range AllCollections {
		collectionName := string(collection)
		if !namesMap[collectionName] {
			fmt.Println("Creating collection", collectionName)
			DB.CreateCollection(Ctx(), collectionName)
		}
	}

	fmt.Println("Database collection check succeeded")
}
