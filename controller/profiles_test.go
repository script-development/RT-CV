package controller

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProfileRoutes(t *testing.T) {
	app := newTestingRouter(t)

	// Get all profiles
	_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/profiles`, TestReqOpts{})

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
		profileRoute := `/api/v1/profiles/` + listProfile.ID.Hex()
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
		_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/profiles`, TestReqOpts{})
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
	_, res = app.MakeRequest(routeBuilder.Post, `/api/v1/profiles`, TestReqOpts{Body: body})
	resProfile := &models.Profile{}
	err = json.Unmarshal(res, resProfile)
	NoError(t, err)
	NotNil(t, resProfile.ID)
	Equal(t, profileToInsert.Name, resProfile.Name)

	// Check if we can fetch the newly inserted profile
	_, res = app.MakeRequest(routeBuilder.Get, `/api/v1/profiles/`+resProfile.ID.Hex(), TestReqOpts{})
	resProfile = &models.Profile{}
	err = json.Unmarshal(res, resProfile)
	NoError(t, err)
	Equal(t, profileToInsert.Name, resProfile.Name)
}

func TestRouteGetProfilesCount(t *testing.T) {
	app := newTestingRouter(t)

	// Get all profiles
	_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/profiles/count`, TestReqOpts{})
	Equal(t, `{"total":2,"usable":1}`, string(res))
}

func TestRouteUpdateProfile(t *testing.T) {
	app := newTestingRouter(t)

	// Get all the profiles from the database
	_, resBody := app.MakeRequest(routeBuilder.Get, `/api/v1/profiles`, TestReqOpts{})
	allProfiles := []models.Profile{}
	err := json.Unmarshal(resBody, &allProfiles)
	NoError(t, err)
	Len(t, allProfiles, 2, "The mock data contains 2 profiles")

	type M map[string]any

	tests := []struct {
		name   string
		req    M
		assert func(t *testing.T, before, after models.Profile)
	}{
		{
			"Nothing should be update on no input",
			M{},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, before.Name, after.Name)
				Equal(t, before.Active, after.Active)
				Equal(t, before.AllowedScrapers, after.AllowedScrapers)
				Equal(t, before.MustDesiredProfession, after.MustDesiredProfession)
				Equal(t, before.DesiredProfessions, after.DesiredProfessions)
				Equal(t, before.YearsSinceWork, after.YearsSinceWork)
				Equal(t, before.MustExpProfession, after.MustExpProfession)
				Equal(t, before.ProfessionExperienced, after.ProfessionExperienced)
				Equal(t, before.MustDriversLicense, after.MustDriversLicense)
				Equal(t, before.DriversLicenses, after.DriversLicenses)
				Equal(t, before.MustEducationFinished, after.MustEducationFinished)
				Equal(t, before.MustEducation, after.MustEducation)
				Equal(t, before.YearsSinceEducation, after.YearsSinceEducation)
				Equal(t, before.Educations, after.Educations)
				Equal(t, before.Zipcodes, after.Zipcodes)
				Equal(t, before.OnMatch, after.OnMatch)
			},
		},
		{
			"Update profile name",
			M{"name": "new name"},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, "new name", after.Name)
			},
		},
		{
			"Set profile to inactive",
			M{"active": false},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, false, after.Active)
			},
		},
		{
			"Set Allowed scraper keys",
			M{"allowedScrapers": []string{mock.Key3.ID.Hex()}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, []primitive.ObjectID{mock.Key3.ID}, after.AllowedScrapers)
			},
		},
		{
			"Set MustDesiredProfession",
			M{"mustDesiredProfession": true},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, true, after.MustDesiredProfession)
			},
		},
		{
			"Set DesiredProfessions",
			M{"desiredProfessions": []models.ProfileProfession{{Name: "updated desired profession"}}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, []models.ProfileProfession{{Name: "updated desired profession"}}, after.DesiredProfessions)
			},
		},
		{
			"Set YearsSinceWork",
			M{"updateYearsSinceWork": M{"yearsSinceWork": 5}},
			func(t *testing.T, before, after models.Profile) {
				yearsSinceWork := 5
				Equal(t, &yearsSinceWork, after.YearsSinceWork)
			},
		},
		{
			"Set MustExpProfession",
			M{"mustExpProfession": true},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, true, after.MustExpProfession)
			},
		},
		{
			"Set ProfessionExperienced",
			M{"professionExperienced": []models.ProfileProfession{{Name: "updated experienced profession"}}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, []models.ProfileProfession{{Name: "updated experienced profession"}}, after.ProfessionExperienced)
			},
		},
		{
			"Set MustDriversLicense",
			M{"mustDriversLicense": true},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, true, after.MustDriversLicense)
			},
		},
		{
			"Set DriversLicenses",
			M{"driversLicenses": []models.ProfileDriversLicense{{Name: "updated drivers license"}}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, []models.ProfileDriversLicense{{Name: "updated drivers license"}}, after.DriversLicenses)
			},
		},
		{
			"Set MustEducationFinished",
			M{"mustEducationFinished": true},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, true, after.MustEducationFinished)
			},
		},
		{
			"Set MustEducation",
			M{"mustEducation": true},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, true, after.MustEducation)
			},
		},
		{
			"Set YearsSinceEducation",
			M{"yearsSinceEducation": 5},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, 5, after.YearsSinceEducation)
			},
		},
		{
			"Set Educations",
			M{"educations": []models.ProfileEducation{{Name: "updated education"}}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, []models.ProfileEducation{{Name: "updated education"}}, after.Educations)
			},
		},
		{
			"Set Zipcodes",
			M{"zipcodes": []models.ProfileDutchZipcode{{From: 1500, To: 2500}}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, []models.ProfileDutchZipcode{{From: 1500, To: 2500}}, after.Zipcodes)
			},
		},
		{
			"Set OnMatch",
			M{"onMatch": models.ProfileOnMatch{
				SendMail: []models.ProfileSendEmailData{{Email: "some-email@f2f.com"}},
			}},
			func(t *testing.T, before, after models.Profile) {
				Equal(t, models.ProfileOnMatch{
					SendMail: []models.ProfileSendEmailData{{Email: "some-email@f2f.com"}},
				}, after.OnMatch)
			},
		},
	}

	for idx, profile := range allProfiles {
		for _, testCase := range tests {
			testCase := testCase
			t.Run(fmt.Sprintf("profile %d %s", idx, testCase.name), func(t *testing.T) {
				reqBodybytes, err := json.Marshal(testCase.req)
				NoError(t, err)

				_, res := app.MakeRequest(routeBuilder.Put, `/api/v1/profiles/`+profile.ID.Hex(), TestReqOpts{
					Body: reqBodybytes,
				})

				var updatedProfile models.Profile
				err = json.Unmarshal(res, &updatedProfile)
				NoError(t, err)

				testCase.assert(t, profile, updatedProfile)
			})
		}
	}
}
