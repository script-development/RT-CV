package controller

import (
	"encoding/json"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestProfileRoutes(t *testing.T) {
	app := newTestingRouter(t)

	// Get all profiles
	_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/control/profiles`, TestReqOpts{})

	// Check if the response contains the profiles inserted in the mock data
	resProfiles := []models.Profile{}
	err := json.Unmarshal(res, &resProfiles)
	NoError(t, err)
	Len(t, resProfiles, 2) // The mock data contains 2 profiles

	// get current profiles in the
	allProfilesInDB := []models.Profile{}
	err = app.db.Find(&models.Profile{}, &allProfilesInDB, nil)
	NoError(t, err)

	// Check if the profiles in the db matches the repsone
	allProfilesInDBJson, err := json.Marshal(allProfilesInDB)
	NoError(t, err)
	Equal(t, string(allProfilesInDBJson), string(res))

	// Get each profile from earlier by id
	for _, listProfile := range resProfiles {
		profileRoute := `/api/v1/control/profiles/` + listProfile.ID.Hex()
		_, res = app.MakeRequest(routeBuilder.Get, profileRoute, TestReqOpts{})

		resProfile := &models.Profile{}
		err = json.Unmarshal(res, resProfile)
		NoError(t, err)
		Equal(t, listProfile.ID.Hex(), resProfile.ID.Hex())

		// Delete the profile and check if it's really deleted
		// Firstly we count how many document we have before the delete
		profilesCountBeforeDeletion := len(resProfiles)

		// Send the delete request
		app.MakeRequest(routeBuilder.Delete, profileRoute, TestReqOpts{})

		// Count how many profiles we have after the deletion
		_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/control/profiles`, TestReqOpts{})
		resProfiles = []models.Profile{}
		err = json.Unmarshal(res, &resProfiles)
		NoError(t, err)

		Equal(t, profilesCountBeforeDeletion-1, len(resProfiles))
	}

	// Try to insert profile
	profileToInsert := models.Profile{
		Name: "newly inserted profile",
		OnMatch: models.ProfileOnMatch{
			HTTPCall: []models.ProfileHTTPCallData{{
				URI:    "http://localhost",
				Method: "GET",
			}},
		},
	}
	body, err := json.Marshal(profileToInsert)
	NoError(t, err)
	_, res = app.MakeRequest(routeBuilder.Post, `/api/v1/control/profiles`, TestReqOpts{Body: body})
	resProfile := &models.Profile{}
	err = json.Unmarshal(res, resProfile)
	NoError(t, err)
	NotNil(t, resProfile.ID)
	Equal(t, profileToInsert.Name, resProfile.Name)

	// Check if we can fetch the newly inserted profile
	_, res = app.MakeRequest(routeBuilder.Get, `/api/v1/control/profiles/`+resProfile.ID.Hex(), TestReqOpts{})
	resProfile = &models.Profile{}
	err = json.Unmarshal(res, resProfile)
	NoError(t, err)
	Equal(t, profileToInsert.Name, resProfile.Name)
}

func TestRouteGetProfilesCount(t *testing.T) {
	app := newTestingRouter(t)

	// Get all profiles
	_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/control/profiles/count`, TestReqOpts{})
	Equal(t, `{"total":2,"usable":2}`, string(res))
}
