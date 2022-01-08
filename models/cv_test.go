package models

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/joho/godotenv"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/emailservice"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func tryLoadEmailEnv() {
	envFileName := ".env"
	_, err := os.Stat(envFileName)
	if err != nil {
		envFileName = "../.env"
		_, err = os.Stat(envFileName)
		if err != nil {
			return
		}
	}

	env, err := godotenv.Read(envFileName)
	if err != nil {
		return
	}

	// Set mail env vars
	for key, value := range env {
		if !strings.HasPrefix(key, "EMAIL_") {
			continue
		}
		if os.Getenv(key) != "" {
			continue
		}
		os.Setenv(key, value)
	}
}

func getExampleCV() *CV {
	now := jsonHelpers.RFC3339Nano(time.Now()).ToPtr()
	return &CV{
		Title:           "Pilot with experience in farming simulator 2020",
		ReferenceNumber: "4455-PIETER",
		CreatedAt:       now,
		LastChanged:     now,

		Educations: []Education{{
			Name:        "Education name",
			Description: "Education description",
			Institute:   "Institute name",
			IsCompleted: true,
			HasDiploma:  true,
			StartDate:   now,
			EndDate:     now,
		}},
		Courses: []Course{{
			Name:        "Course name",
			Institute:   "Institute name",
			IsCompleted: true,
			StartDate:   now,
			EndDate:     now,
		}},
		WorkExperiences: []WorkExperience{{
			Description:       "WorkExperience description",
			Profession:        "hitman",
			StartDate:         now,
			EndDate:           now,
			StillEmployed:     true,
			Employer:          "Bond.. James bond",
			WeeklyHoursWorked: 60,
		}},
		PreferredJobs: []string{"hitman"},
		Languages: []Language{{
			Name:         "Language name",
			LevelSpoken:  LanguageLevelExcellent,
			LevelWritten: LanguageLevelGood,
		}},
		Competences: []Competence{{
			Name:        "Competence name",
			Description: "Competence description",
		}},
		Interests: []Interest{{
			Name:        "Interest name",
			Description: "Interest description",
		}},
		PersonalPresentation: "Sir",
		DriversLicenses: []jsonHelpers.DriversLicense{
			jsonHelpers.NewDriversLicense("AAA"),
		},
		PersonalDetails: PersonalDetails{
			Initials:          "P.S.",
			FirstName:         "D.R. Pietter",
			SurNamePrefix:     "Ven ther",
			SurName:           "Steen",
			DateOfBirth:       now,
			Gender:            "Apache helicopter",
			StreetName:        "Streetname abc",
			HouseNumber:       "33",
			HouseNumberSuffix: "b",
			Zip:               "9999AB",
			City:              "Groningen",
			Country:           "Netherlands",
			PhoneNumber:       &jsonHelpers.PhoneNumber{IsLocal: true, Number: 611223344},
			Email:             "dr.p.steen@smart-people.com",
		},
	}
}

func TestGetEmailHTML(t *testing.T) {
	matchTest := "this is a test text that should re-appear in the response html"

	cv := getExampleCV()

	profileObjectID := primitive.NewObjectID()
	profile := Profile{
		M:       db.M{ID: profileObjectID},
		Name:    "profile name",
		Domains: []string{"test.com"},
	}

	htmlBuff, err := cv.GetEmailHTML(profile, matchTest)
	NoError(t, err)

	html := htmlBuff.String()
	Contains(t, html, matchTest)
	Contains(t, html, cv.PersonalDetails.FirstName+" "+cv.PersonalDetails.SurName)
	Contains(t, html, cv.PersonalDetails.Email)
	Contains(t, html, cv.PersonalDetails.PhoneNumber.String())
	Contains(t, html, profile.Name)
	Contains(t, html, cv.ReferenceNumber)
	Contains(t, html, profile.ID.Hex())
}

func TestGetEmailAttachmentHTML(t *testing.T) {
	cv := getExampleCV()

	htmlBuff, err := cv.GetEmailAttachmentHTML()
	NoError(t, err)

	html := htmlBuff.String()

	Contains(t, html, cv.FullName())
	Contains(t, html, "Referentie: #"+cv.ReferenceNumber)
	Contains(t, html, "Laatst gewijzigd: ")
	Contains(t, html, cv.PersonalDetails.Gender)
	Contains(t, html, cv.PersonalDetails.StreetName+" "+cv.PersonalDetails.HouseNumber+" "+cv.PersonalDetails.HouseNumberSuffix)
	Contains(t, html, cv.PersonalDetails.Zip+" "+cv.PersonalDetails.City)
	Contains(t, html, cv.PersonalDetails.Email)
	Contains(t, html, cv.PersonalDetails.PhoneNumber.String())
	Contains(t, html, cv.PersonalPresentation)

	Contains(t, html, cv.WorkExperiences[0].Profession)
	Contains(t, html, cv.WorkExperiences[0].Employer)
	Contains(t, html, " - heden")
	Contains(t, html, cv.WorkExperiences[0].Description)

	Contains(t, html, cv.Educations[0].Institute)
	Contains(t, html, cv.Educations[0].Description)

	Contains(t, html, cv.Courses[0].Institute)
	Contains(t, html, cv.Courses[0].Description)

	Contains(t, html, "<th>Taal</th>")
	Contains(t, html, "<th>Mondeling</th>")
	Contains(t, html, "<th>Schriftelijk</th>")
	Contains(t, html, cv.Languages[0].Name)
	Contains(t, html, cv.Languages[0].LevelSpoken.String())
	Contains(t, html, cv.Languages[0].LevelWritten.String())

	Contains(t, html, cv.DriversLicenses[0].String())

	Contains(t, html, cv.Interests[0].Name)
	Contains(t, html, cv.Interests[0].Description)
}

func getBaseProjectPath() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	if strings.HasSuffix(p, "/models") || strings.HasSuffix(p, "\\models") {
		p = path.Clean(path.Join(p, ".."))
	}
	return p
}

func TestGetNewEmailAttachmentPDF(t *testing.T) {
	pdfGeneratorProjectPath := path.Join(getBaseProjectPath(), "pdf_generator")

	pdfGeneratorBin := path.Join(pdfGeneratorProjectPath, "bin/pdf_generator.exe")
	pdfOutfile := path.Join(pdfGeneratorProjectPath, "example.pdf")

	_, err := os.Open(pdfGeneratorBin)
	if os.IsNotExist(err) {
		t.Skip(pdfGeneratorBin + " does not exist, skipping test")
	} else if err != nil {
		NoError(t, err)
	}

	cv := getExampleCV()
	jsonCV, err := json.Marshal(cv)
	NoError(t, err)

	tests := [][]string{
		{"--dummy"},
		{"--data", string(jsonCV)},
	}

	for _, args := range tests {
		os.Remove(pdfOutfile)

		cmd := exec.Command(pdfGeneratorBin, append(args, "--out", pdfOutfile)...)
		cmd.Dir = pdfGeneratorProjectPath
		out, err := cmd.CombinedOutput()
		NoError(t, err, string(out))
		NotEmpty(t, out)

		_, err = os.Open(pdfOutfile)
		NoError(t, err)
	}
}

func TestSendMail(t *testing.T) {
	tryLoadEmailEnv()

	emailConf := emailservice.EmailServerConfigurationFromEnv()
	if emailConf.Host == "" || emailConf.From == "" {
		t.Skip("Missing email server env variables to test sending emails")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Initialize the mail service
	err := emailservice.Setup(
		emailConf,
		func(err error) {
			NoError(t, err)
			wg.Done()
		},
	)
	if err != nil {
		log.WithError(err).Error("Error initializing email service")
		return
	}

	cv := getExampleCV()
	profile := Profile{
		M:       db.M{ID: primitive.NewObjectID()},
		Name:    "profile name",
		Domains: []string{"test.com"},
	}

	emailBody, err := cv.GetEmailHTML(profile, "on data from the void")
	NoError(t, err)

	emailToSendData := &ProfileSendEmailData{Email: "example@localhost"}
	err = emailToSendData.SendEmail(profile, emailBody.Bytes(), nil)
	NoError(t, err)

	// Wait for the email to succeed
	wg.Wait()
}
