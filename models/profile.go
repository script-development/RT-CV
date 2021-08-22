package models

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Profile contains all the information about a search profile
type Profile struct {
	db.M        `bson:"inline"`
	Name        string   `json:"name"`
	Active      bool     `json:"active"`
	Domains     []string `json:"domains"`
	ListProfile bool     `json:"-" bson:"-"` // TODO find out what this is

	MustDesiredProfession bool                `json:"mustDesiredProfession"`
	DesiredProfessions    []ProfileProfession `json:"desiredProfessions"`

	YearsSinceWork        *int `json:"yearsSinceWork"`
	MustExpProfession     bool `json:"mustExpProfession"`
	ProfessionExperienced []ProfileProfession

	MustDriversLicense bool `json:"mustDriversLicense"`
	DriversLicenses    []ProfileDriversLicense

	MustEducationFinished bool               `json:"mustEducationFinished"`
	MustEducation         bool               `json:"mustEducation"`
	YearsSinceEducation   int                `json:"yearsSinceEducation"`
	Educations            []ProfileEducation `json:"educations"`

	Emails   []ProfileEmail        `json:"emails"`
	Zipcodes []ProfileDutchZipcode `json:"zipCodes"`
}

// CollectionName returns the collection name of the Profile
func (*Profile) CollectionName() string {
	return "profiles"
}

// DefaultFindFilters returns the default filters for the Find function
func (*Profile) DefaultFindFilters() bson.M {
	return bson.M{
		"active": true,
	}
}

// GetProfiles returns all profiles from the database
func GetProfiles(conn db.Connection) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, nil)
	return profiles, err
}

// ProfileProfession contains information about a proffession
type ProfileProfession struct {
	Name string `json:"name"`

	// TODO find out what this is about?
	// HeadFunctionID int
	// SubsectorLevel1ID int
	// SubsectorLevel2ID int
	// SubsectorLevel3ID int
	// SubsectorLevel4ID int
	// SubsectorLevel5ID int
	// SubsectorLevel6ID int
}

// ProfileDriversLicense contains the drivers license name
type ProfileDriversLicense struct {
	Name string `json:"name"`
}

// ProfileEducation contains information about an education
type ProfileEducation struct {
	Name string `json:"name"`
	// HeadEducationID int
	// SubsectorID     int
}

// ProfileEmail only contains an email address
type ProfileEmail struct {
	Email string `json:"email"`
}

// type ProfileProfession struct {
// 	ID        int `gorm:"primaryKey"`
// 	ProfileID int
// 	Name      string
// }

// ProfileDutchZipcode is dutch zipcode range limited to the number
type ProfileDutchZipcode struct {
	From uint16 `json:"from"`
	To   uint16 `json:"to"`
}
