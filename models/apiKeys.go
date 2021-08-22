package models

import (
	"encoding/json"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

// APIKey contains a registered API key
type APIKey struct {
	db.M    `bson:"inline"`
	Enabled bool       `json:"enabled"`
	Domains []string   `json:"domains"`
	Key     string     `json:"key"`
	Roles   APIKeyRole `json:"roles"`
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

// APIKeyRole is a role that tells what someone can and can't do
// Roles can be combined together using bit sifting
// For example:
//   APIKeyRoleScraper | APIKeyRoleInformationObtainer // Is a valid APIKeyRole that represends 2 rules
type APIKeyRole uint64

const (
	// APIKeyRoleScraper can insert scraped data
	// = 1
	APIKeyRoleScraper APIKeyRole = 1 << iota

	// APIKeyRoleInformationObtainer can obtain information the server has
	// = 2
	APIKeyRoleInformationObtainer

	// APIKeyRoleController can control the server
	// = 4
	APIKeyRoleController

	// APIKeyRoleAdmin Currently unused
	// = 8
	APIKeyRoleAdmin
)

var (
	// APIKeyRoleAll contains all of the above roles and thus can access everything
	APIKeyRoleAll = APIKeyRoleScraper | APIKeyRoleInformationObtainer | APIKeyRoleController | APIKeyRoleAdmin
	// APIKeyRoleAllArray is an array of all roles
	APIKeyRoleAllArray = []APIKeyRole{APIKeyRoleScraper, APIKeyRoleInformationObtainer, APIKeyRoleController, APIKeyRoleAdmin}
)

// Description returns a description of the role
// Only works on single roles
func (a APIKeyRole) Description() (description string, ok bool) {
	switch a {
	case APIKeyRoleScraper:
		return "Can insert scraped data", true
	case APIKeyRoleInformationObtainer:
		return "Can obtain information the server has", true
	case APIKeyRoleController:
		return "Can control the server", true
	case APIKeyRoleAdmin:
		return "Unused role", true
	default:
		return "Unknown role", false
	}
}

// APIRole contains information about a APIKeyRole
type APIRole struct {
	Role        APIKeyRole `json:"role"`
	Description string     `json:"description"`
}

// MarshalJSON convers the unreadable role number into an array of APIRole
func (a APIKeyRole) MarshalJSON() ([]byte, error) {
	res := []APIRole{}
	for _, role := range APIKeyRoleAllArray {
		if a&role == role {
			description, _ := role.Description()
			res = append(res, APIRole{Role: role, Description: description})
		}
	}

	if len(res) == 0 {
		return []byte(`[]`), nil
	}
	return json.Marshal(res)
}

// ContainsAll check if a contains all of other
func (a APIKeyRole) ContainsAll(other APIKeyRole) bool {
	return a&other == other
}

// ContainsSome check if a contains some of other
func (a APIKeyRole) ContainsSome(other APIKeyRole) bool {
	return a&other > 0
}
