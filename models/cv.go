package models

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/helpers/schema"
)

// CV contains all information that belongs to a curriculum vitae
// TODO check the json removed fields if we actually should use them
type CV struct {
	Title                string                   `json:"-"`
	ReferenceNumber      string                   `json:"-"`
	CreatedAt            *jsonHelpers.RFC3339Nano `json:"-"`
	LastChanged          *jsonHelpers.RFC3339Nano `json:"-"`
	Educations           []Education              `json:"educations"`
	Courses              []Course                 `json:"courses"`
	WorkExperiences      []WorkExperience         `json:"workExperiences"`
	PreferredJobs        []string                 `json:"preferredJobs"`
	Languages            []Language               `json:"languages"`
	Competences          []Competence             `json:"-"`
	Interests            []Interest               `json:"-"`
	PersonalDetails      PersonalDetails          `json:"personalDetails"`
	PersonalPresentation string                   `json:"-"`
	DriversLicenses      []string                 `json:"driversLicenses"`
}

// Education is something a user has followed
type Education struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// TODO find difference between isCompleted and hasdiploma
	IsCompleted bool                     `json:"isCompleted"`
	HasDiploma  bool                     `json:"hasDiploma"`
	StartDate   *jsonHelpers.RFC3339Nano `json:"startDate"`
	EndDate     *jsonHelpers.RFC3339Nano `json:"endDate"`
	Institute   string                   `json:"institute"`
	SubsectorID int                      `json:"subsectorID"`
}

// Course is something a user has followed
type Course struct {
	Name           string                   `json:"name"`
	NormalizedName string                   `json:"normalizedName"`
	StartDate      *jsonHelpers.RFC3339Nano `json:"startDate"`
	EndDate        *jsonHelpers.RFC3339Nano `json:"endDate"`
	IsCompleted    bool                     `json:"isCompleted"`
	Institute      string                   `json:"institute"`
	Description    string                   `json:"description"`
}

// WorkExperience is experience in work
type WorkExperience struct {
	Description       string                   `json:"description"`
	Profession        string                   `json:"profession"`
	StartDate         *jsonHelpers.RFC3339Nano `json:"startDate"`
	EndDate           *jsonHelpers.RFC3339Nano `json:"endDate"`
	StillEmployed     bool                     `json:"stillEmployed"`
	Employer          string                   `json:"employer"`
	WeeklyHoursWorked int                      `json:"weeklyHoursWorked"`
}

// LanguageLevel is something that i'm not sure what it is
type LanguageLevel uint

// The lanague levels available
const (
	LanguageLevelUnknown LanguageLevel = iota
	LanguageLevelReasonable
	LanguageLevelGood
	LanguageLevelExcellent
)

const langLevelDescription = `0. Unknown
1. Reasonable
2. Good
3. Excellent`

// Valid returns weather the language level is valid
func (l LanguageLevel) Valid() bool {
	return l >= LanguageLevelUnknown && l <= LanguageLevelExcellent
}

// JSONSchemaDescribe implements schema.Describe
func (LanguageLevel) JSONSchemaDescribe() schema.Property {
	return schema.Property{
		Title:       "Language level",
		Description: langLevelDescription,
		Type:        schema.PropertyTypeNumber,
		Enum: []interface{}{
			LanguageLevelUnknown,
			LanguageLevelReasonable,
			LanguageLevelGood,
			LanguageLevelExcellent,
		},
	}
}

// Language is a language a user can speak
type Language struct {
	Name         string        `json:"name"`
	LevelSpoken  LanguageLevel `json:"levelSpoken"`
	LevelWritten LanguageLevel `json:"levelWritten"`
}

// Competence is an activity a user is "good" at
type Competence struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Interest contains a job the user is interested in
type Interest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PersonalDetails contains personal info
type PersonalDetails struct {
	Initials          string                   `json:"initials" jsonSchema:"notRequired"`
	FirstName         string                   `json:"firstName"`
	SurNamePrefix     string                   `json:"surNamePrefix" jsonSchema:"notRequired"`
	SurName           string                   `json:"surName" jsonSchema:"notRequired"`
	DateOfBirth       *jsonHelpers.RFC3339Nano `json:"dob" jsonSchema:"notRequired"`
	Gender            string                   `json:"gender" jsonSchema:"notRequired"`
	StreetName        string                   `json:"streetName" jsonSchema:"notRequired"`
	HouseNumber       string                   `json:"houseNumber" jsonSchema:"notRequired"`
	HouseNumberSuffix string                   `json:"houseNumberSuffix" jsonSchema:"notRequired"`
	Zip               string                   `json:"zip" jsonSchema:"notRequired"`
	City              string                   `json:"city" jsonSchema:"notRequired"`
	Country           string                   `json:"country" jsonSchema:"notRequired"`
	PhoneNumber       string                   `json:"phoneNumber" jsonSchema:"notRequired"`
	Email             string                   `json:"email" jsonSchema:"notRequired"`
}

// GetHTML generates a HTML document from the input cv
func (cv *CV) GetHTML(profile Profile, matchText string) (*bytes.Buffer, error) {
	tmpl, err := template.ParseFiles("./assets/email-template.html")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		// For testing perposes
		tmpl, err = template.ParseFiles("../assets/email-template.html")
		if err != nil {
			return nil, err
		}
	}

	domains := "onbekend"
	if len(profile.Domains) > 0 {
		domains = strings.Join(profile.Domains, ", ")
	}

	input := struct {
		Profile   Profile
		Cv        *CV
		MatchText string
		LogoURL   string
		Domains   string

		// The normal `Profile.ID.String()`` is more of a debug value than a real id value so we add the hex to this field
		ProfileIDHex string
	}{
		Profile:      profile,
		ProfileIDHex: profile.ID.Hex(),
		Cv:           cv,
		MatchText:    matchText,
		LogoURL:      os.Getenv("EMAIL_LOGO_URL"),
		Domains:      domains,
	}

	buff := bytes.NewBuffer(nil)
	err = tmpl.Execute(buff, input)
	return buff, err
}

// GetPDF generates a PDF from a cv that can be send
func (cv *CV) GetPDF(profile Profile, matchText string) ([]byte, error) {
	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	generator.MarginBottom.Set(50)
	generator.MarginTop.Set(0)
	generator.MarginLeft.Set(0)
	generator.MarginRight.Set(0)
	generator.ImageQuality.Set(100)

	html, err := cv.GetHTML(profile, matchText)
	if err != nil {
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(html)
	page.PageOptions = wkhtmltopdf.NewPageOptions()
	page.DisableSmartShrinking.Set(true)
	generator.AddPage(page)

	err = generator.Create()
	if err != nil {
		return nil, err
	}

	return generator.Bytes(), nil
}

// Validate validates the cv and returns an error if it's not valid
func (cv *CV) Validate() error {
	now := time.Now()
	if cv.CreatedAt != nil && cv.CreatedAt.Time().After(now) {
		return errors.New("createdAt can't be in the future")
	}
	if cv.LastChanged != nil && cv.LastChanged.Time().After(now) {
		return errors.New("lastChanged can't be in the future")
	}
	if cv.PersonalDetails.DateOfBirth != nil && cv.PersonalDetails.DateOfBirth.Time().After(now) {
		return errors.New("dateOfBirth can't be in the future")
	}

	for idx, lang := range cv.Languages {
		if !lang.LevelSpoken.Valid() {
			return fmt.Errorf("languages.%d.levelSpoken is invalid", idx)
		}
		if !lang.LevelWritten.Valid() {
			return fmt.Errorf("languages.%d.levelWritten is invalid", idx)
		}
	}

	return nil
}
