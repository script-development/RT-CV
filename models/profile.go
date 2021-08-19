package models

import (
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
)

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
	Zipcodes              []Zipcode
}

func (*Profile) CollectionName() string {
	return "profiles"
}
func (*Profile) DefaultFindFilters() bson.M {
	return bson.M{
		"active": true,
	}
}

func GetProfiles(conn db.Connection) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, nil)
	return profiles, err
}

type Profession struct {
	Name           string
	HeadFunctionID int

	// TODO find out what this is about?
	// SubsectorLevel1ID int
	// SubsectorLevel2ID int
	// SubsectorLevel3ID int
	// SubsectorLevel4ID int
	// SubsectorLevel5ID int
	// SubsectorLevel6ID int
}

type DriversLicense struct {
	Name string
}

type DBEducation struct {
	Name string
	// HeadEducationID int
	// SubsectorID     int
}

type Email struct {
	Email string
}

// type ProfileProfession struct {
// 	ID        int `gorm:"primaryKey"`
// 	ProfileID int
// 	Name      string
// }

type Zipcode struct {
	From uint16
	To   uint16
}

// var now = time.Now()
// var mockGetProfiles = []Profile{
// 	{
// 		ID:                    primitive.NewObjectID(),
// 		Name:                  "Mock profile 1",
// 		YearsSinceWork:        nil,
// 		Active:                true,
// 		MustExpProfession:     true,
// 		MustDesiredProfession: false,
// 		MustEducation:         true,
// 		MustEducationFinished: true,
// 		MustDriversLicense:    true,
// 		Domains:               []string{"werk.nl"},
// 		ListProfile:           true,
// 		YearsSinceEducation:   1,
// 		DesiredProfessions: []Profession{{
// 			Name: "Rapper",
// 		}},
// 		ProfessionExperienced: []Profession{{
// 			Name: "Dancer",
// 		}},
// 		DriversLicenses: []DriversLicense{{
// 			Name: "A",
// 		}},
// 		Educations: []DBEducation{{
// 			Name: "Default",
// 		}},
// 		Emails: []Email{{
// 			Email: "abc@example.com",
// 		}},
// 		Zipcodes: []Zipcode{{
// 			From: 2000,
// 			To:   8000,
// 		}},
// 	},
// 	{
// 		ID:                    primitive.NewObjectID(),
// 		Name:                  "Mock profile 2",
// 		YearsSinceWork:        nil,
// 		Active:                true,
// 		MustExpProfession:     false,
// 		MustDesiredProfession: false,
// 		MustEducation:         false,
// 		MustEducationFinished: false,
// 		MustDriversLicense:    false,
// 		Domains:               []string{"werk.nl"},
// 		ListProfile:           false,
// 		YearsSinceEducation:   0,
// 		DesiredProfessions:    nil,
// 		ProfessionExperienced: nil,
// 		DriversLicenses:       nil,
// 		Educations:            nil,
// 		Emails:                nil,
// 		Zipcodes:              nil,
// 	},
// }
