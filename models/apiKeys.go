package models

import (
	"github.com/script-development/RT-CV/db"
)

type ApiKey struct {
	ID      int `gorm:"primaryKey"`
	Enabled bool
	SiteId  uint
	Site    Site
	Key     string
	Roles   ApiKeyRole
}

func (ApiKey) TableName() string {
	return "api_keys"
}

func GetApiKeys() ([]ApiKey, error) {
	keys := []ApiKey{}
	err := db.DB.
		Preload("Site").
		Where("enabled = 1").
		Find(&keys).Error

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
