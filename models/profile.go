package models

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/script-development/RT-CV/db"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
)

// Profile contains all the information about a search profile
type Profile struct {
	db.M        `bson:",inline"`
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

	Zipcodes []ProfileDutchZipcode `json:"zipCodes"`

	// What should happen on a match
	OnMatch ProfileOnMatch `json:"onMatch"`
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

// ProfileOnMatch defines what should happen when a profile is matched to a CV
type ProfileOnMatch struct {
	SendMail []ProfileSendEmailData `json:"sendMail"`
	HTTPCall []ProfileHTTPCallData  `json:"httpCall"`
}

// ProfileSendEmailData only contains an email address atm
type ProfileSendEmailData struct {
	Email string `json:"email"`
}

// SendEmail sends an email
func (*ProfileSendEmailData) SendEmail() {
	// FIXME implment this
}

// ProfileHTTPCallData defines a http address that should be called when a match was made
type ProfileHTTPCallData struct {
	URI    string `json:"uri"`
	Method string `json:"method"`
}

// MakeRequest creates a http request
func (d *ProfileHTTPCallData) MakeRequest() {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(d.URI)
	req.Header.SetMethod(d.Method)

	// FIXME set request timeout
	// FIXME set body data or url data

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	fasthttp.Do(req, resp)
}

// ValidateCreateNewProfile validates a new profile to create
func (p *Profile) ValidateCreateNewProfile() error {
	// TODO this needs more validation

	if p.Name == "" {
		return errors.New("name must be set")
	}

	if len(p.OnMatch.SendMail) == 0 && len(p.OnMatch.HTTPCall) == 0 {
		return errors.New("at least on of the profile onMatch options be set")
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
	for idx, call := range p.OnMatch.HTTPCall {
		uri, err := url.Parse(call.URI)
		if err != nil {
			return fmt.Errorf("onMatch.httpCall[%d].uri: %s", idx, err.Error())
		}
		if uri.Scheme != "http" && uri.Scheme != "https" {
			return fmt.Errorf("onMatch.httpCall[%d].uri: url schema must be set to http or https", idx)
		}
		if uri.User != nil {
			return fmt.Errorf("onMatch.httpCall[%d].uri: user information is not allowed", idx)
		}
		if uri.Host == "" && uri.Opaque == "" {
			return fmt.Errorf("onMatch.httpCall[%d].uri: host must be set", idx)
		}
		switch call.Method {
		case "", "GET", "POST", "PATCH", "PUT", "DELETE":
		default:
			return fmt.Errorf(
				"onMatch.httpCall[%d].method: not a valid method, must be one of "+
					`"GET", "POST", "PATCH", "PUT", "DELETE" or empty to default to GET`,
				idx,
			)
		}
	}

	return nil
}
