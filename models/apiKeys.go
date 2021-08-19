package models

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

type ApiKey struct {
	db.M    `bson:"inline"`
	Enabled bool
	Domains []string
	Key     string
	Roles   ApiKeyRole
}

func (a *ApiKey) CollectionName() string {
	return "apiKeys"
}

func (m *ApiKey) DefaultFindFilters() bson.M {
	return bson.M{
		"enabled": true,
	}
}

func GetApiKeys(conn db.Connection) ([]ApiKey, error) {
	keys := []ApiKey{}
	err := conn.Find(&ApiKey{}, &keys, nil)
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

// var mockGetApiKeys = []ApiKey{
// 	{
// 		M:       db.NewM(),
// 		Enabled: true,
// 		Domains: []string{"werk.nl"},
// 		Key:     "abc",
// 		Roles:   ApiKeyRoleScraper,
// 	},
// }
