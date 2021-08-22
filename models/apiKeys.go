package models

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

// APIKey contains a registered API key
type APIKey struct {
	db.M    `bson:"inline"`
	Enabled bool
	Domains []string
	Key     string
	Roles   APIKeyRole
}

// CollectionName returns the collection name of the ApiKey
func (*APIKey) CollectionName() string {
	return "apiKeys"
}

// DefaultFindFilters sets the default filters for the db connection find
func (*APIKey) DefaultFindFilters() bson.M {
	return bson.M{
		"enabled": true,
	}
}

// GetAPIKeys returns all the keys registered in the database
func GetAPIKeys(conn db.Connection) ([]APIKey, error) {
	keys := []APIKey{}
	err := conn.Find(&APIKey{}, &keys, nil)
	return keys, err
}

// APIKeyRole is a role
type APIKeyRole uint64

const (
	// APIKeyRoleScraper can access the scraper routes
	APIKeyRoleScraper APIKeyRole = 1 << iota // 1

	// APIKeyRoleInformationObtainer can obtain information the server has
	APIKeyRoleInformationObtainer // 2

	// APIKeyRoleController can control server settings
	APIKeyRoleController // 4

	// APIKeyRoleAdmin can do administrative tasks
	APIKeyRoleAdmin // 8
)

var (
	// APIKeyRoleAll contains all of the above roles and thus can access everything
	APIKeyRoleAll = APIKeyRoleScraper | APIKeyRoleInformationObtainer | APIKeyRoleController | APIKeyRoleAdmin
)

// APIRole contains information about a APIKeyRole
type APIRole struct {
	Role        APIKeyRole `json:"role"`
	Description string     `json:"description"`
}

// APIRoles contains all api roles with a description
var APIRoles = []APIRole{
	{
		APIKeyRoleScraper,
		"Can insert scraped data",
	},
	{
		APIKeyRoleInformationObtainer,
		"Can obtain scraped information",
	},
	{
		APIKeyRoleController,
		"Can obtain scraped information",
	},
	{
		APIKeyRoleAdmin,
		"Admin (Currently unused)",
	},
}

// ContainsAll check if a contains all of other
func (a APIKeyRole) ContainsAll(other APIKeyRole) bool {
	return a&other == other
}

// ContainsSome check if a contains some of other
func (a APIKeyRole) ContainsSome(other APIKeyRole) bool {
	return a&other > 0
}
