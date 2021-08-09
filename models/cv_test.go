package models

import (
	"strconv"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
)

func TestGetHtml(t *testing.T) {
	matchTest := "this is a test text that should re-appear in the response html"

	cv := Cv{
		Title:           "Pilot with experience in farming simulator 2020",
		ReferenceNumber: "4455-PIETER",

		PersonalDetails: PersonalDetails{
			Initials:          "P.S.",
			FirstName:         "D.R. Pietter",
			SurNamePrefix:     "Ven ther",
			SurName:           "Steen",
			DateOfBirth:       time.Now().Format(time.RFC3339Nano),
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

	profile := Profile{
		ID:   223344,
		Name: "profile name",
		Site: Site{
			Domain: "test.com",
		},
	}

	htmlBuff, err := cv.GetHtml(profile, matchTest)
	if err != nil {
		NoError(t, err)
		return
	}

	html := htmlBuff.String()
	Contains(t, html, matchTest)
	Contains(t, html, cv.PersonalDetails.FirstName+" "+cv.PersonalDetails.SurName)
	Contains(t, html, cv.PersonalDetails.Email)
	Contains(t, html, cv.PersonalDetails.PhoneNumber)
	Contains(t, html, profile.Name)
	Contains(t, html, cv.ReferenceNumber)
	Contains(t, html, strconv.Itoa(profile.ID))
}
