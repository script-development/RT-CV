package models

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ScraperLoginUsers defines all the users a scraper can use
type ScraperLoginUsers struct {
	db.M      `bson:",inline"`
	ScraperID primitive.ObjectID `json:"scraperId" bson:"scraperId"`
	Users     []ScraperLoginUser `json:"users"`
}

// CollectionName should yield the collection name for the entry
func (*ScraperLoginUsers) CollectionName() string {
	return "scraperLoginUsers"
}

// Indexes implements db.Entry
func (*ScraperLoginUsers) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{Keys: bson.M{"scraperId": 1}},
	}
}

// ScraperLoginUser defines a user that can be used by a scraper to login into a scraped website
type ScraperLoginUser struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}
