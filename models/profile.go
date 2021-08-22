package models

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Profile contains all the information about a search profile
type Profile struct {
	db.M                  `bson:"inline"`
	Name                  string
	YearsSinceWork        *int
	Active                bool
	MustExpProfession     bool
	MustDesiredProfession bool
	MustEducation         bool
	MustEducationFinished bool
	MustDriversLicense    bool
	Domains               []string
	ListProfile           bool // TODO find out what this is
	YearsSinceEducation   int

	DesiredProfessions    []ProfileProfession
	ProfessionExperienced []ProfileProfession
	DriversLicenses       []ProfileDriversLicense
	Educations            []ProfileEducation
	Emails                []ProfileEmail
	Zipcodes              []ProfileDutchZipcode
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
	Name string

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
	Name string
}

// ProfileEducation contains information about an education
type ProfileEducation struct {
	Name string
	// HeadEducationID int
	// SubsectorID     int
}

// ProfileEmail only contains an email address
type ProfileEmail struct {
	Email string
}

// type ProfileProfession struct {
// 	ID        int `gorm:"primaryKey"`
// 	ProfileID int
// 	Name      string
// }

// ProfileDutchZipcode is dutch zipcode range limited to the number
type ProfileDutchZipcode struct {
	From uint16
	To   uint16
}
