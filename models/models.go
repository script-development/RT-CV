package models

// Toggle this to true in a test to use mock data
var Testing = false

type Site struct {
	ID     uint `gorm:"primaryKey"`
	Domain string
}

func (Site) TableName() string {
	return "sites"
}

type User struct {
	ID           int `gorm:"primaryKey"`
	Username     string
	Password     string `gorm:"column:PASSWORD"`
	Active       bool
	Session      string `gorm:"column:SESSION"`
	InUse        bool   `gorm:"column:in_use"`
	Consumer     string
	LastModified interface{} `gorm:"column:last_modified;type:time"`
	SiteID       int         `gorm:"column:site_id"`
}

func (User) TableName() string {
	return "accounts"
}
