package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		log.Fatal(err)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	ctxCancel()
	if err != nil {
		log.Fatal(err)
	}

	ctx, ctxCancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	ctxCancel()
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(os.Getenv("MONGODB_DATABASE"))
	fmt.Println("Connected to database")
}
