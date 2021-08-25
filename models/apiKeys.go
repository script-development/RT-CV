package models

import (
	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIKey contains a registered API key
type APIKey struct {
	db.M    `bson:",inline"`
	Enabled bool       `json:"enabled"`
	Domains []string   `json:"domains"`
	Key     string     `json:"key"`
	Roles   APIKeyRole `json:"roles"`

	// System indicates if this is a key required by the system
	// These are keys whereof at least one needs to exists otherwise RT-CV would not work
	System bool `json:"system"`
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

// APIKeyInfo contains information about an API key
// This key can be send to someone safely without exposing the key
// To generate this object use the (*APIKey).Info method
type APIKeyInfo struct {
	ID      primitive.ObjectID `json:"id"`
	Domains []string           `json:"domains"`
	Roles   []APIRole          `json:"roles"`
	System  bool               `json:"system"`
}

// Info converts the APIKey into APIKeyInfo
// This key can be send to someone safely without exposing the key
func (a *APIKey) Info() APIKeyInfo {
	return APIKeyInfo{
		ID:      a.ID,
		Domains: a.Domains,
		Roles:   a.Roles.ConvertToAPIRoles(),
		System:  a.System,
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

	// APIKeyRoleDashboard can access the dashboard and modify server state
	// = 8
	APIKeyRoleDashboard

	// APIKeyRoleAdmin Currently unused
	// = 16
	APIKeyRoleAdmin
)

var (
	// APIKeyRoleAll contains all of the above roles and thus can access everything
	APIKeyRoleAll = APIKeyRoleScraper | APIKeyRoleInformationObtainer | APIKeyRoleController | APIKeyRoleDashboard | APIKeyRoleAdmin
	// APIKeyRoleAllArray is an array of all roles
	APIKeyRoleAllArray = []APIKeyRole{
		APIKeyRoleScraper,
		APIKeyRoleInformationObtainer,
		APIKeyRoleController,
		APIKeyRoleDashboard,
		APIKeyRoleAdmin,
	}
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
	case APIKeyRoleDashboard:
		return "Can access the dashboard and modify server state", true
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

// ConvertToAPIRoles convers the unreadable role number into an array of APIRole
func (a APIKeyRole) ConvertToAPIRoles() []APIRole {
	res := []APIRole{}
	for _, role := range APIKeyRoleAllArray {
		if a&role == role {
			description, _ := role.Description()
			res = append(res, APIRole{Role: role, Description: description})
		}
	}
	return res
}

// ContainsAll check if a contains all of other
func (a APIKeyRole) ContainsAll(other APIKeyRole) bool {
	return a&other == other
}

// ContainsSome check if a contains some of other
func (a APIKeyRole) ContainsSome(other APIKeyRole) bool {
	return a&other > 0
}

// Valid returns if the role is valid role
// Empty roles are also invalid
func (a APIKeyRole) Valid() bool {
	return a > 0 && a <= APIKeyRoleAll
}

// CheckDashboardKeyExists checks weather the required system keys are available and if not creates them
func CheckDashboardKeyExists(conn db.Connection) {
	keys := []APIKey{}
	err := conn.Find(&APIKey{}, &keys, bson.M{"system": true, "roles": APIKeyRoleDashboard})
	if err != nil {
		log.WithError(err).Fatalf("unable to fetch api keys")
	}

	if len(keys) != 0 {
		log.Infof("One system dashboard key exists with id %s and role %d", keys[0].ID.Hex(), APIKeyRoleDashboard)
		return
	}

	log.Info("System dashboard key does not yet exists, creating one..")
	key := &APIKey{
		M:       db.NewM(),
		Enabled: true,
		Domains: []string{"*"},
		Key:     string(random.GenerateKey()),
		Roles:   APIKeyRoleDashboard,
		System:  true,
	}
	err = conn.Insert(key)
	if err != nil {
		log.WithError(err).Fatalf("Unable to insert dashboard system api keys")
	}
	log.WithField("key", key.Key).WithField("id", key.ID.Hex()).Info("Created system key")
}
