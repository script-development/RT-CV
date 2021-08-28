package models

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
)

// Period is a period of time from a date to another date
type Period struct {
	Start   jsonHelpers.RFC3339Nano `json:"start"`
	End     jsonHelpers.RFC3339Nano `json:"end"`
	Present bool                    `json:"present"`
}

// CV contains all information that belongs to a curriculum vitae
// TODO check the json removed fields if we actually should use them
type CV struct {
	Title                string                   `json:"-"`
	ReferenceNumber      string                   `json:"-"`
	CreatedAt            *jsonHelpers.RFC3339Nano `json:"-"`
	LastChanged          jsonHelpers.RFC3339Nano  `json:"-"`
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
	Period      Period                   `json:"period"`
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
	Period            Period                   `json:"period"`
	StartDate         *jsonHelpers.RFC3339Nano `json:"startDate"`
	EndDate           *jsonHelpers.RFC3339Nano `json:"endDate"`
	StillEmployed     bool                     `json:"stillEmployed"`
	Employer          string                   `json:"employer"`
	WeeklyHoursWorked int                      `json:"weeklyHoursWorked"`
}

// LanguageProficiency is something that i'm not sure what it is
// FIXME
type LanguageProficiency int

// Language is a language a user can speak
type Language struct {
	Name         string              `json:"name"`
	LevelSpoken  LanguageProficiency `json:"levelSpoken"`
	LevelWritten LanguageProficiency `json:"levelWritten"`
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
	Initials          string                  `json:"initials" jsonSchema:"notRequired"`
	FirstName         string                  `json:"firstName"`
	SurNamePrefix     string                  `json:"surNamePrefix" jsonSchema:"notRequired"`
	SurName           string                  `json:"surName" jsonSchema:"notRequired"`
	DateOfBirth       jsonHelpers.RFC3339Nano `json:"dob" jsonSchema:"notRequired"`
	Gender            string                  `json:"gender" jsonSchema:"notRequired"`
	StreetName        string                  `json:"streetName" jsonSchema:"notRequired"`
	HouseNumber       string                  `json:"houseNumber" jsonSchema:"notRequired"`
	HouseNumberSuffix string                  `json:"houseNumberSuffix" jsonSchema:"notRequired"`
	Zip               string                  `json:"zip" jsonSchema:"notRequired"`
	City              string                  `json:"city" jsonSchema:"notRequired"`
	Country           string                  `json:"country" jsonSchema:"notRequired"`
	PhoneNumber       string                  `json:"phoneNumber" jsonSchema:"notRequired"`
	Email             string                  `json:"email" jsonSchema:"notRequired"`
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
		LogoURL:      os.Getenv("LOGO"),
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
