package models

import (
	"context"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApiKeysCollection() *mongo.Collection {
	return db.DB.Collection("apiKeys")
}

type ApiKey struct {
	ID      primitive.ObjectID `bson:"_id"`
	Enabled bool
	Domains []string
	Key     string
	Roles   ApiKeyRole
}

func GetApiKeys() ([]ApiKey, error) {
	if Testing {
		panic("FIXME")
	}

	c, err := ApiKeysCollection().Find(db.Ctx(), bson.M{
		"enabled": true,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close(context.Background())

	keys := []ApiKey{}
	err = c.All(db.Ctx(), &keys)
	return keys, err
}

type ApiKeyRole uint64

const (
	ApiKeyRoleScraper             = 1 << iota // 1
	ApiKeyRoleInformationObtainer             // 2
	ApiKeyRoleController                      // 4
	ApiKeyRoleAdmin                           // 8
)

type ApiRole struct {
	Role        ApiKeyRole `json:"role"`
	Description string     `json:"description"`
}

var ApiRoles = []ApiRole{
	{
		ApiKeyRoleScraper,
		"Can insert scraped data",
	},
	{
		ApiKeyRoleInformationObtainer,
		"Can obtain scraped information",
	},
	{
		ApiKeyRoleController,
		"Can obtain scraped information",
	},
	{
		ApiKeyRoleAdmin,
		"Admin (Currently unused)",
	},
}

func (a ApiKeyRole) ContainsAll(other ApiKeyRole) bool {
	return a&other == other
}

func (a ApiKeyRole) ContainsSome(other ApiKeyRole) bool {
	return a&other > 0
}
