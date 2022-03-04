package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/jordan-wright/email"
	fuzzymatcher "github.com/mjarkk/fuzzy-matcher"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/emailservice"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"jaytaylor.com/html2text"
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
	YearsSinceEducation   int                `json:"yearsSinceEducation" bson:"yearsSinceEducation"`
	Educations            []ProfileEducation `json:"educations" bson:"educations"`

	Zipcodes []ProfileDutchZipcode `json:"zipCodes" bson:"zipCodes"`

	// What should happen on a match
	OnMatch ProfileOnMatch `json:"onMatch" bson:"onMatch" description:"What should happen when a match is made on this profile"`

	// OldID is used to keep track of converted old profiles
	OldID *uint64 `bson:"_old_id" json:"-"`

	// Variables set by the matching process only when they needed
	// These are mainly used for caching so we don't have to calculate values twice
	// There values where detected using the -profile flag, see main.go for more info
	EducationFuzzyMatcher             *fuzzymatcher.Matcher        `bson:"-" json:"-"`
	ProfessionExperiencedFuzzyMatcher *fuzzymatcher.Matcher        `bson:"-" json:"-"`
	DesiredProfessionsFuzzyMatcher    *fuzzymatcher.Matcher        `bson:"-" json:"-"`
	DomainPartsCache                  [][]string                   `bson:"-" json:"-"`
	NormalizedDriversLicensesCache    []jsonHelpers.DriversLicense `bson:"-" json:"-"`
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
		{Keys: bson.M{"onMatch.httpCall": 1}},
	}
}

func actualActiveProfilesFilter() bson.M {
	isArrayWContent := bson.M{"$not": bson.M{"$size": 0}, "$type": "array"}
	return bson.M{
		"active": true,
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"desiredProfessions": isArrayWContent},
					{"professionExperienced": isArrayWContent},
					{"driversLicenses": isArrayWContent},
					{"educations": isArrayWContent},
				},
			},
			{
				"$or": []bson.M{
					{"onMatch.sendMail": isArrayWContent},
					{"onMatch.httpCall": isArrayWContent},
				},
			},
		},
	}
}

// GetActualActiveProfiles returns that we can actually use
// Matches are not really helpfull if no desiredProfessions, professionExperienced, driversLicenses or educations is set
// Matches without an onMatch property are useless as we can't send the match anywhere
func GetActualActiveProfiles(conn db.Connection) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, actualActiveProfilesFilter())
	return profiles, err
}

// GetActualActiveProfilesCount does the same as GetActualActiveProfiles but only returns the number of found profiles
func GetActualActiveProfilesCount(conn db.Connection) (uint64, error) {
	return conn.Count(&Profile{}, actualActiveProfilesFilter())
}

// GetProfiles returns all profiles from the database
func GetProfiles(conn db.Connection) ([]Profile, error) {
	profiles := []Profile{}
	err := conn.Find(&Profile{}, &profiles, nil)
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
	SendMail   []ProfileSendEmailData `json:"sendMail" bson:"sendMail"`
	HTTPCall   []ProfileHTTPCallData  `json:"httpCall" bson:"httpCall"`
	PdfOptions *PdfOptions            `json:"pdfOptions" bson:"pdfOptions" description:"Options for customizing the PDF, for more info about this object look at the /tryPdfGenerator page"`
}

// HasPDFOptions returns true if the PdfOptions is set and has at least one option set
func (onMatch *ProfileOnMatch) HasPDFOptions() bool {
	options := onMatch.PdfOptions
	if options == nil {
		return false
	}

	for _, entry := range []*string{
		options.FontHeader,
		options.FontRegular,
		options.Style,
		options.HeaderColor,
		options.SubHeaderColor,
		options.LogoImageURL,
		options.CompanyName,
		options.CompanyAddress,
	} {
		if entry != nil {
			return true
		}
	}

	return false
}

// PdfOptions contains options for the creation of the pdf
type PdfOptions struct {
	// See pdf_generator/bin/fonts.dart > _fontFilesMap for available fonts
	FontHeader  *string `json:"fontHeader" bson:"fontHeader"`
	FontRegular *string `json:"fontRegular" bson:"fontRegular"`

	// See pdf_generator/bin/args.dart > LayoutStyle for available styles
	Style *string `json:"style" bson:"style"`

	// Expected to be a hex value like: #ffffff
	HeaderColor    *string `json:"headerColor" bson:"headerColor"`
	SubHeaderColor *string `json:"subHeaderColor" bson:"subHeaderColor"`

	LogoImageURL   *string `json:"logoImageUrl" bson:"logoImageUrl"`
	CompanyName    *string `json:"companyName" bson:"companyName"`
	CompanyAddress *string `json:"companyAddress" bson:"companyAddress"`
}

// ProfileSendEmailData only contains an email address atm
type ProfileSendEmailData struct {
	Email string `json:"email"`
}

// SendEmail sends an email
func (d *ProfileSendEmailData) SendEmail(profile Profile, htmlBody []byte, pdfFile *os.File) error {
	e := email.NewEmail()

	e.To = []string{d.Email}
	e.Subject = "Nieuwe match voor " + profile.Name
	e.HTML = htmlBody
	text, _ := html2text.FromString(string(htmlBody), html2text.Options{})
	e.Text = []byte(text)

	if pdfFile != nil {
		pdfFile.Seek(0, 0)
		_, err := e.Attach(pdfFile, "match.pdf", "application/pdf")
		if err != nil {
			return err
		}
	}

	emailservice.SendMail(e)
	return nil
}

// ProfileHTTPCallData defines a http address that should be called when a match was made
type ProfileHTTPCallData struct {
	URI    string `json:"uri"`
	Method string `json:"method"`
}

// MakeRequest creates a http request
func (d *ProfileHTTPCallData) MakeRequest(profile Profile, match Match) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(d.URI)
	req.Header.SetMethod(d.Method)

	// FIXME set request timeout
	// FIXME url data in case of get request
	value, err := json.Marshal(map[string]interface{}{
		"profileId": profile.ID.Hex(),
		"match":     match,
	})
	if err != nil {
		req.ResetBody()
		req.AppendBody(value)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// We don't care about the error returned here as it's not our task to keep the endpoint running
	var _ = fasthttp.Do(req, resp)
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
