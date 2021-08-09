package models

import (
	"bytes"
	"html/template"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type Period struct {
	Start   string `json:"start"` // iso 8601 time
	End     string `json:"end"`   // iso 8601 time
	Present bool   `json:"present"`
}

type Cv struct {
	Title                string           `json:"title"`
	ReferenceNumber      string           `json:"referenceNumber"`
	LastChanged          string           `json:"lastChanged"` // iso 8601 time
	Educations           []Education      `json:"educations"`
	Courses              []Course         `json:"courses"`
	WorkExperiences      []WorkExperience `json:"workExperiences"`
	PreferredJobs        []string         `json:"preferredJobs"`
	Languages            []Language       `json:"languages"`
	Competences          []Competence     `json:"competences"`
	Interests            []Interest       `json:"interests"`
	PersonalDetails      PersonalDetails  `json:"personalDetails"`
	PersonalPresentation string           `json:"personalPresentation"`
	DriversLicenses      []string         `json:"driversLicenses"`
	CreatedAt            *string          `json:"createdAt"` // iso 8601 time
}

type Education struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// TODO find difference between iscompleted and hasdiploma
	IsCompleted bool   `json:"isCompleted"`
	HasDiploma  bool   `json:"hasDiploma"`
	Period      Period `json:"period"`
	StartDate   string `json:"startDate"` // iso 8601 time
	EndDate     string `json:"endDate"`   // iso 8601 time
	Institute   string `json:"institute"`
	SubsectorID int    `json:"subsectorID"`
}

type Course struct {
	Name           string `json:"name"`
	NormalizedName string `json:"normalizedName"`
	StartDate      string `json:"startDate"` // iso 8601 time
	EndDate        string `json:"endDate"`   // iso 8601 time
	IsCompleted    bool   `json:"isCompleted"`
	Institute      string `json:"institute"`
	Description    string `json:"description"`
}

type WorkExperience struct {
	Description       string `json:"description"`
	Profession        string `json:"profession"`
	Period            Period `json:"period"`
	StartDate         string `json:"startDate"` // iso 8601 time
	EndDate           string `json:"endDate"`   // iso 8601 time
	StillEmployed     bool   `json:"stillEmployed"`
	Employer          string `json:"employer"`
	WeeklyHoursWorked int    `json:"weeklyHoursWorked"`
}

type LanguageProficiency int

type Language struct {
	Name         string              `json:"name"`
	LevelSpoken  LanguageProficiency `json:"levelSpoken"`
	LevelWritten LanguageProficiency `json:"levelWritten"`
}

type Competence struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Interest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PersonalDetails struct {
	Initials          string `json:"initials"`
	FirstName         string `json:"firstName"`
	SurNamePrefix     string `json:"surNamePrefix"`
	SurName           string `json:"surName"`
	DateOfBirth       string `json:"dob"` // iso 8601 time
	Gender            string `json:"gender"`
	StreetName        string `json:"streetName"`
	HouseNumber       string `json:"houseNumber"`
	HouseNumberSuffix string `json:"houseNumberSuffix"`
	Zip               string `json:"zip"`
	City              string `json:"city"`
	Country           string `json:"country"`
	PhoneNumber       string `json:"phoneNumber"`
	Email             string `json:"email"`
}

// GetHtml generates a HTML document from the input cv
func (cv *Cv) GetHtml(profile Profile, matchText string) (*bytes.Buffer, error) {
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

	input := struct {
		Profile   Profile
		Cv        *Cv
		MatchText string
		LogoUrl   string
	}{
		Profile:   profile,
		Cv:        cv,
		MatchText: matchText,
		LogoUrl:   os.Getenv("LOGO"),
	}

	buff := bytes.NewBuffer(nil)
	err = tmpl.Execute(buff, input)
	return buff, err
}

// GetPDF generates a PDF from a cv that can be send
func (cv *Cv) GetPDF(profile Profile, matchText string) ([]byte, error) {
	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	generator.MarginBottom.Set(50)
	generator.MarginTop.Set(0)
	generator.MarginLeft.Set(0)
	generator.MarginRight.Set(0)
	generator.ImageQuality.Set(100)

	html, err := cv.GetHtml(profile, matchText)
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
