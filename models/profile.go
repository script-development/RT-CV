package models

import (
	"errors"
	"fmt"
	"regexp"

	fuzzymatcher "github.com/mjarkk/fuzzy-matcher"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Profile contains all the information about a search profile
type Profile struct {
	db.M            `bson:",inline"`
	Name            string               `json:"name"`
	Active          bool                 `json:"active"`
	AllowedScrapers []primitive.ObjectID `json:"allowedScrapers" bson:"allowedScrapers" description:"Define a list of scraper keys that can use this profile, if value is undefined or empty all keys are allowed"`

	MustDesiredProfession bool                `json:"mustDesiredProfession" bson:"mustDesiredProfession"`
	DesiredProfessions    []ProfileProfession `json:"desiredProfessions" bson:"desiredProfessions"`

	YearsSinceWork        *int                `json:"yearsSinceWork" bson:"yearsSinceWork"`
	MustExpProfession     bool                `json:"mustExpProfession" bson:"mustExpProfession"`
	ProfessionExperienced []ProfileProfession `json:"professionExperienced" bson:"professionExperienced"`

	MustDriversLicense bool                    `json:"mustDriversLicense" bson:"mustDriversLicense"`
	DriversLicenses    []ProfileDriversLicense `json:"driversLicenses" bson:"driversLicenses"`

	MustEducationFinished bool               `json:"mustEducationFinished" bson:"mustEducationFinished"`
	MustEducation         bool               `json:"mustEducation" bson:"mustEducation" description:"Should a found CV at least have one education regardless of if it's complete"`
	YearsSinceEducation   *int               `json:"yearsSinceEducation" bson:"yearsSinceEducation"`
	Educations            []ProfileEducation `json:"educations" bson:"educations"`

	Zipcodes []ProfileDutchZipcode `json:"zipCodes" bson:"zipCodes"`

	// What should happen on a match
	OnMatch ProfileOnMatch `json:"onMatch" bson:"onMatch" description:"What should happen when a match is made on this profile"`

	Lables map[string]any `json:"labels" description:"custom labels that can be used by API users to identify profiles, the key needs to be a string and the value can be anything"`

	ListsAllowed bool `json:"listsAllowed" bson:"listsAllowed"`

	// OldID is used to keep track of converted old profiles
	OldID *uint64 `bson:"_old_id" json:"-"`

	// Variables set by the matching process only when they needed
	// These are mainly used for caching so we don't have to calculate values twice
	// There values where detected using the -profile flag, see main.go for more info
	EducationFuzzyMatcherCache             *fuzzymatcher.Matcher        `bson:"-" json:"-"`
	ProfessionExperiencedFuzzyMatcherCache *fuzzymatcher.Matcher        `bson:"-" json:"-"`
	DesiredProfessionsFuzzyMatcherCache    *fuzzymatcher.Matcher        `bson:"-" json:"-"`
	DomainPartsCache                       [][]string                   `bson:"-" json:"-"`
	NormalizedDriversLicensesCache         []jsonHelpers.DriversLicense `bson:"-" json:"-"`

	// Tell if this profile should use complex search
}

// CollectionName returns the collection name of the Profile
func (*Profile) CollectionName() string {
	return "profiles"
}

// Indexes implements db.Entry
func (*Profile) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{Keys: bson.M{"active": 1}},
		{Keys: bson.M{"desiredProfessions": 1}},
		{Keys: bson.M{"professionExperienced": 1}},
		{Keys: bson.M{"driversLicenses": 1}},
		{Keys: bson.M{"educations": 1}},
		{Keys: bson.M{"zipCodes": 1}},
		{Keys: bson.M{"onMatch.sendMail": 1}},
		{Keys: bson.M{"listsAllowed": 1}},
	}
}

var isArrayWContent = bson.M{"$not": bson.M{"$size": 0}, "$type": "array"}

// GetListsProfiles returns all profiles that can be used for the cv lists functionality
func GetListsProfiles(conn db.Connection) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, bson.M{
		"active":       true,
		"listsAllowed": true,
		"zipCodes":     isArrayWContent,
	})
	return profiles, err
}

func actualActiveMatchProfilesFilter() bson.M {
	return bson.M{
		"active": true,
		"$or": []bson.M{
			{"desiredProfessions": isArrayWContent},
			{"professionExperienced": isArrayWContent},
			{"driversLicenses": isArrayWContent},
			{"educations": isArrayWContent},
		},
		// we use $not here as there are properties without this property and with `$not: true` we match `false`, `undefined` and `null
		"listsAllowed": bson.M{"$not": bson.M{"$eq": true}},
	}
}

// GetActualMatchActiveProfiles returns that we can actually use
// Matches are not really helpfull if no desiredProfessions, professionExperienced, driversLicenses or educations is set
// Matches without an onMatch property are useless as we can't send the match anywhere
func GetActualMatchActiveProfiles(conn db.Connection) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, actualActiveMatchProfilesFilter())
	return profiles, err
}

// GetActualMatchActiveProfilesCount does the same as GetActualMatchActiveProfiles but only returns the number of found profiles
func GetActualMatchActiveProfilesCount(conn db.Connection) (uint64, error) {
	return conn.Count(&Profile{}, actualActiveMatchProfilesFilter())
}

// GetProfiles returns all profiles from the database
func GetProfiles(conn db.Connection, filters primitive.M) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, filters)
	return profiles, err
}

// GetProfilesCount returns the amount of profiles in the database
func GetProfilesCount(conn db.Connection) (uint64, error) {
	return conn.Count(&Profile{}, nil)
}

// GetProfile returns a profile by id
func GetProfile(conn db.Connection, id primitive.ObjectID) (Profile, error) {
	profile := Profile{}
	err := conn.FindOne(&profile, bson.M{"_id": id})
	return profile, err
}

// ProfileProfession contains information about a proffession
type ProfileProfession struct {
	Name   string `json:"name"`
	LeafId primitive.ObjectID
}

// ProfileDriversLicense contains the drivers license name
type ProfileDriversLicense struct {
	Name string `json:"name"`
}

// ProfileEducation contains information about an education
type ProfileEducation struct {
	Name string `json:"name"`
}

// ProfileDutchZipcode is dutch zipcode range limited to the number
type ProfileDutchZipcode struct {
	From uint16 `json:"from"`
	To   uint16 `json:"to"`
}

// IsWithinCithAndArea checks if the cityAndArea provided are in the range range of the zipcode
func (p *ProfileDutchZipcode) IsWithinCithAndArea(cityAndArea uint16) bool {
	if p.From > p.To {
		// Swap from and to
		p.From, p.To = p.To, p.From
	}

	if cityAndArea < 1_000 || cityAndArea >= 10_000 {
		// Illegal postal code
		return false
	}
	return p.From <= cityAndArea && p.To >= cityAndArea
}

// ProfileOnMatch defines what should happen when a profile is matched to a CV
type ProfileOnMatch struct {
	SendMail []ProfileSendEmailData `json:"sendMail" bson:"sendMail"`
}

// ProfileSendEmailData only contains an email address atm
type ProfileSendEmailData struct {
	Email string `json:"email"`
}

// CheckAPIKeysExists checks if apiKeys are valid IDs of existing keys
func CheckAPIKeysExists(conn db.Connection, apiKeys []primitive.ObjectID) error {
	if len(apiKeys) == 0 {
		return nil
	}

	apiKeysInDB, err := GetAPIKeys(conn)
	if err != nil {
		return err
	}
outer:
	for _, allowedKey := range apiKeys {
		for _, apiKey := range apiKeysInDB {
			if allowedKey == apiKey.ID {
				continue outer
			}
		}
		return fmt.Errorf("unknown api key id %s", allowedKey.Hex())
	}
	return nil
}

// ValidateCreateNewProfile validates a new profile to create
func (p *Profile) ValidateCreateNewProfile(conn db.Connection) error {
	// TODO this needs more validation

	if p.Name == "" {
		return errors.New("name must be set")
	}

	if len(p.AllowedScrapers) > 0 {
		err := CheckAPIKeysExists(conn, p.AllowedScrapers)
		if err != nil {
			return err
		}
	}

	emailRegex := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@" +
			"[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
			"(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
	for idx, mail := range p.OnMatch.SendMail {
		if len(mail.Email) < 3 || len(mail.Email) > 254 || !emailRegex.MatchString(mail.Email) {
			return fmt.Errorf("onMatch.sendMail[%d].email: invalid email address", idx)
		}
	}

	return nil
}
