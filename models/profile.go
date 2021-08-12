package models

import (
	"time"

	"github.com/script-development/RT-CV/db"
)

type Profile struct {
	ID                    int `gorm:"primaryKey"`
	Name                  string
	CompanyID             int
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

	DesiredProfessions    []ProfileDesiredProfession
	ProfessionExperienced []ProfileProfessionExperience
	DriversLicenses       []ProfileDriversLicense
	Educations            []ProfileEducation
	Emails                []ProfileEmail
	Zipcodes              []ProfileZipcode
}

func (Profile) TableName() string {
	return "profiles"
}

func GetProfiles() ([]Profile, error) {
	if Testing {
		return []Profile{}, nil
	}

	profiles := []Profile{}
	err := db.DB.
		Preload("DesiredProfessions.Profession").
		Preload("ProfessionExperienced.Profession").
		Preload("DriversLicenses.DriversLicense").
		Preload("Educations.Education").
		Preload("Emails.Email").
		Preload("Zipcodes.Zipcode").
		Preload("Site").
		Where("active = 1 AND deleted_at IS NULL").
		Find(&profiles).Error
	return profiles, err
}

type ProfileDesiredProfession struct {
	ID           int `gorm:"primaryKey"`
	ProfileID    int
	ProfessionID int
	Profession   Profession
}

func (ProfileDesiredProfession) TableName() string {
	return "profile_desired_profession"
}

type Profession struct {
	ID             int `gorm:"primaryKey"`
	Name           string
	HeadFunctionID int

	// TODO find out what this is about?
	SubsectorLevel1ID int
	SubsectorLevel2ID int
	SubsectorLevel3ID int
	SubsectorLevel4ID int
	SubsectorLevel5ID int
	SubsectorLevel6ID int
}

func (Profession) TableName() string {
	return "professions"
}

type ProfileProfessionExperience struct {
	ID           int `gorm:"primaryKey"`
	ProfileID    int
	ProfessionID int
	Profession   Profession
}

func (ProfileProfessionExperience) TableName() string {
	return "profile_profession_experience"
}

type ProfileDriversLicense struct {
	ID               int `gorm:"primaryKey"`
	ProfileID        int
	DriversLicenseID int
	DriversLicense   DriversLicense
}

func (ProfileDriversLicense) TableName() string {
	return "profile_drivers_license"
}

type DriversLicense struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func (DriversLicense) TableName() string {
	return "drivers_licenses"
}

type ProfileEducation struct {
	ID          int `gorm:"primaryKey"`
	ProfileID   int
	EducationID int
	Education   DBEducation
}

func (ProfileEducation) TableName() string {
	return "profile_education"
}

type DBEducation struct {
	ID              int `gorm:"primaryKey"`
	Name            string
	HeadEducationID int
	SubsectorID     int
}

func (DBEducation) TableName() string {
	return "educations"
}

type ProfileEmail struct {
	ID        int `gorm:"primaryKey"`
	ProfileID int
	EmailID   int
	Email     Email
}

func (ProfileEmail) TableName() string {
	return "profile_email"
}

type Email struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	CompanyID string
}

func (Email) TableName() string {
	return "emails"
}

// type ProfileProfession struct {
// 	ID        int `gorm:"primaryKey"`
// 	ProfileID int
// 	Name      string
// }

type ProfileZipcode struct {
	ID        int `gorm:"primaryKey"`
	ProfileID int
	ZipcodeID int
	Zipcode   Zipcode
}

func (ProfileZipcode) TableName() string {
	return "profile_zipcode"
}

type Zipcode struct {
	ID   int `gorm:"primaryKey"`
	From int
	To   int
}

func (Zipcode) TableName() string {
	return "zipcodes"
}
