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

	DesiredProfessions    []Profession
	ProfessionExperienced []Profession
	DriversLicenses       []DriversLicense
	Educations            []DBEducation
	Emails                []Email
	Zipcodes              []DutchZipcode
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

// Profession contains information about a proffession
type Profession struct {
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

// DriversLicense contains the drivers license name
type DriversLicense struct {
	Name string
}

// DBEducation contains education name
type DBEducation struct {
	Name string
	// HeadEducationID int
	// SubsectorID     int
}

// Email only contains an email address
type Email struct {
	Email string
}

// type ProfileProfession struct {
// 	ID        int `gorm:"primaryKey"`
// 	ProfileID int
// 	Name      string
// }

// DutchZipcode is dutch zipcode range limited to the number
type DutchZipcode struct {
	From uint16
	To   uint16
}
