package models

import (
	"testing"
	"time"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetEmailHTML(t *testing.T) {
	matchTest := "this is a test text that should re-appear in the response html"

	cv := CV{
		Title:           "Pilot with experience in farming simulator 2020",
		ReferenceNumber: "4455-PIETER",

		PersonalDetails: PersonalDetails{
			Initials:          "P.S.",
			FirstName:         "D.R. Pietter",
			SurNamePrefix:     "Ven ther",
			SurName:           "Steen",
			DateOfBirth:       jsonHelpers.RFC3339Nano(time.Now()).ToPtr(),
			Gender:            "Apache helicopter",
			StreetName:        "Streetname abc",
			HouseNumber:       "33",
			HouseNumberSuffix: "b",
			Zip:               "9999AB",
			City:              "Groningen",
			Country:           "Netherlands",
			PhoneNumber:       "06-11223344",
			Email:             "dr.p.steen@smart-people.com",
		},
	}

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
	Contains(t, html, cv.PersonalDetails.PhoneNumber)
	Contains(t, html, profile.Name)
	Contains(t, html, cv.ReferenceNumber)
	Contains(t, html, profile.ID.Hex())
}

func TestGetEmailAttachmentHTML(t *testing.T) {
	now := jsonHelpers.RFC3339Nano(time.Now()).ToPtr()

	cv := CV{
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
		DriversLicenses:      []string{"AAA"},

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
			PhoneNumber:       "06-11223344",
			Email:             "dr.p.steen@smart-people.com",
		},
	}

	profileObjectID := primitive.NewObjectID()
	profile := Profile{
		M:       db.M{ID: profileObjectID},
		Name:    "profile name",
		Domains: []string{"test.com"},
	}

	htmlBuff, err := cv.GetEmailAttachmentHTML(profile)
	NoError(t, err)

	html := htmlBuff.String()

	Contains(t, html, cv.FullName())
	Contains(t, html, "Referentie: #"+cv.ReferenceNumber)
	Contains(t, html, "Laatst gewijzigd: ")
	Contains(t, html, cv.PersonalDetails.Gender)
	Contains(t, html, cv.PersonalDetails.StreetName+" "+cv.PersonalDetails.HouseNumber+" "+cv.PersonalDetails.HouseNumberSuffix)
	Contains(t, html, cv.PersonalDetails.Zip+" "+cv.PersonalDetails.City)
	Contains(t, html, cv.PersonalDetails.Email)
	Contains(t, html, cv.PersonalDetails.PhoneNumber)
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

	Contains(t, html, cv.DriversLicenses[0])

	Contains(t, html, cv.Interests[0].Name)
	Contains(t, html, cv.Interests[0].Description)
}
