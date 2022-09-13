package models

import (
	"testing"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/testingdb"
	. "github.com/tj/assert"
)

func testProfilesSetupDB(t *testing.T) *testingdb.TestConnection {
	testingDB := testingdb.NewDB()

	zipCodes := []ProfileDutchZipcode{{1000, 1999}}
	desiredProfessions := []ProfileProfession{{Name: "gangster"}}

	err := testingDB.UnsafeInsert(
		&Profile{ // active list profile
			M:                  db.NewM(),
			Active:             true,
			Zipcodes:           zipCodes,
			DesiredProfessions: desiredProfessions,
			ListsAllowed:       false,
		},
		&Profile{ // IN-active list profile
			M:                  db.NewM(),
			Active:             false,
			Zipcodes:           zipCodes,
			DesiredProfessions: desiredProfessions,
			ListsAllowed:       false,
		},
		&Profile{ // active list profile
			M:                  db.NewM(),
			Active:             true,
			Zipcodes:           zipCodes,
			DesiredProfessions: desiredProfessions,
			ListsAllowed:       true,
		},
		&Profile{ // IN-active list profile
			M:                  db.NewM(),
			Active:             false,
			Zipcodes:           zipCodes,
			DesiredProfessions: desiredProfessions,
			ListsAllowed:       true,
		},
	)
	NoError(t, err)

	return testingDB
}

func TestGetActualMatchActiveProfiles(t *testing.T) {
	d := testProfilesSetupDB(t)
	profiles, err := GetActualMatchActiveProfiles(d)
	NoError(t, err)
	Len(t, profiles, 1)
}

func TestGetListsProfiles(t *testing.T) {
	d := testProfilesSetupDB(t)
	profiles, err := GetListsProfiles(d)
	NoError(t, err)
	Len(t, profiles, 1)
}
