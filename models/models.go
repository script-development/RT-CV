package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Toggle this to true in a test to use mock data
var Testing = false

type Site struct {
	ID     primitive.ObjectID `bson:"_id"`
	Domain string
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string
	Password     string
	Active       bool
	Session      string
	InUse        bool
	Consumer     string
	LastModified *time.Time
	SiteID       int
}
