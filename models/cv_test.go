package models

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/apex/log"
	"github.com/joho/godotenv"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/emailservice"
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

func TestGetEmailHTML(t *testing.T) {
	matchTest := "this is a test text that should re-appear in the response html"

	cv := ExampleCV()

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
	cv := ExampleCV()

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

	_, err := os.Open(pdfGeneratorBin)
	if os.IsNotExist(err) {
		t.Skip(pdfGeneratorBin + " does not exist, skipping test")
	} else if err != nil {
		NoError(t, err)
	}

	cv := ExampleCV()

	ptrStr := func(in string) *string { return &in }

	optionsToTest := []*PdfOptions{
		nil,
		{},
		{
			FontHeader:  ptrStr("IBMPlexSerif"),
			FontRegular: ptrStr("IBMPlexSerif"),
		},
		{
			FontHeader: ptrStr("This font does not exist, pdf generator should use fallback font and not fail"),
		},
		{
			Style: ptrStr("style_2"),
		},
		{
			Style: ptrStr("This style does not exist, pdf generator should use fallback style and not fail"),
		},
		{
			HeaderColor:    ptrStr("#FFFFFF"),
			SubHeaderColor: ptrStr("#FFF"),
		},
		{
			HeaderColor:    ptrStr("#ffffff"),
			SubHeaderColor: ptrStr("#fff"),
		},
		{
			HeaderColor:    ptrStr("FFFFFF"),
			SubHeaderColor: ptrStr("FFF"),
		},
		{
			HeaderColor:    ptrStr("ffffff"),
			SubHeaderColor: ptrStr("fff"),
		},
		{
			CompanyName: ptrStr("A company name"),
		},
		{
			CompanyAddress: ptrStr("A company address"),
		},
	}

	for _, options := range optionsToTest {
		jsonOptionsBytes, err := json.Marshal(options)
		NoError(t, err)
		jsonOptions := string(jsonOptionsBytes)

		file, err := cv.GetPDF(options, &pdfGeneratorProjectPath)
		if file != nil {
			_, err = os.Open(file.Name())
			os.Remove(file.Name())
			NoError(t, err, jsonOptions)
		}
		NoError(t, err, jsonOptions)
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

	cv := ExampleCV()
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
