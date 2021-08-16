package models

import (
	"time"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProfilesCollection() *mongo.Collection {
	return db.DB.Collection("profiles")
}

type Profile struct {
	ID                    primitive.ObjectID `bson:"_id"`
	Name                  string
	YearsSinceWork        *int
	Active                bool
	MustExpProfession     bool
	MustDesiredProfession bool
	MustEducation         bool
	MustEducationFinished bool
	MustDriversLicense    bool
	CreatedAt             *time.Time
	UpdatedAt             *time.Time
	DeletedAt             *time.Time
	SiteId                *int
	Site                  Site
	ListProfile           bool
	YearsSinceEducation   int

	DesiredProfessions    []Profession
	ProfessionExperienced []Profession
	DriversLicenses       []DriversLicense
	Educations            []DBEducation
	Emails                []Email
	Zipcodes              []Zipcode
}

func GetProfiles() ([]Profile, error) {
	if Testing {
		panic("FIXME")
	}

	c, err := ProfilesCollection().Find(db.Ctx(), bson.M{
		"active": true,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close(db.Ctx())

	profiles := []Profile{}
	err = c.All(db.Ctx(), &profiles)
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
	Name string
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
