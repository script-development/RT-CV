package mock

// Mock provides mock a testing db with mock data
//
// The mock data should be predictable

import (
	"time"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/testingdb"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// Key1 is a mock api key
	Key1 = &models.APIKey{
		M:       db.NewM(),
		Name:    "Key with all roles",
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "aaa",
		Roles:   models.APIKeyRoleAll,
	}
	// Key2 is a mock api key
	Key2 = &models.APIKey{
		M:       db.NewM(),
		Name:    "Scraper key",
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "bbb",
		Roles:   models.APIKeyRoleScraper,
	}
	// Key3 is a mock api key
	Key3 = &models.APIKey{
		M:       db.NewM(),
		Name:    "Information obtainer key",
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "ccc",
		Roles:   models.APIKeyRoleInformationObtainer,
	}
	// DashboardKey is the mock key for the dashboard
	DashboardKey = &models.APIKey{
		M: db.M{
			ID: primitive.ObjectID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
		Name:    "System login key",
		System:  true,
		Enabled: true,
		Domains: []string{"*"},
		Key:     "ddd",
		Roles:   models.APIKeyRoleDashboard,
	}

	// Profile1 contains the first example profile
	yearsSinceEducation := 1
	Profile1 = &models.Profile{
		M:                     db.NewM(),
		Name:                  "Mock profile 1",
		AllowedScrapers:       []primitive.ObjectID{Key2.ID},
		YearsSinceWork:        nil,
		Active:                true,
		MustExpProfession:     true,
		MustDesiredProfession: false,
		MustEducation:         true,
		MustEducationFinished: true,
		MustDriversLicense:    true,
		YearsSinceEducation:   &yearsSinceEducation,
		DesiredProfessions: []models.ProfileProfession{{
			Name: "Rapper",
		}},
		ProfessionExperienced: []models.ProfileProfession{{
			Name: "Dancer",
		}},
		DriversLicenses: []models.ProfileDriversLicense{{
			Name: "A",
		}},
		Educations: []models.ProfileEducation{{
			Name: "Default",
		}},
		OnMatch: models.ProfileOnMatch{
			SendMail: []models.ProfileSendEmailData{
				{Email: "example@script.nl"},
			},
		},
		Zipcodes: []models.ProfileDutchZipcode{{
			From: 2000,
			To:   8000,
		}},
	}

	// Profile2 contains the second example profile
	Profile2 = &models.Profile{
		M:                     db.NewM(),
		Name:                  "Mock profile 2",
		AllowedScrapers:       []primitive.ObjectID{Key2.ID},
		YearsSinceWork:        nil,
		Active:                true,
		MustExpProfession:     false,
		MustDesiredProfession: false,
		MustEducation:         false,
		MustEducationFinished: false,
		MustDriversLicense:    false,
		YearsSinceEducation:   nil,
		DesiredProfessions:    nil,
		ProfessionExperienced: nil,
		DriversLicenses:       nil,
		Educations:            nil,
		Zipcodes:              nil,
		OnMatch: models.ProfileOnMatch{
			SendMail: []models.ProfileSendEmailData{
				{Email: "example@script.nl"},
			},
		},
	}

	// mockMatch1 contains a example match between profile 1 and a cv
	mockMatch1 = &models.Match{
		M:                     db.NewM(),
		RequestID:             primitive.NewObjectID(),
		ProfileID:             Profile1.ID,
		KeyID:                 Key1.ID,
		When:                  jsonHelpers.RFC3339Nano(time.Now().Add(-(time.Minute * 15))),
		ReferenceNr:           "a",
		ProfessionExperienced: &professionExperienced,
		YearsSinceEducation:   &yearsSinceEducation,
		Education:             &matchedEducation,
		DriversLicense:        true,
	}
	mockMatch2 = &models.Match{
		M:           db.NewM(),
		RequestID:   primitive.NewObjectID(),
		ProfileID:   Profile2.ID,
		KeyID:       Key2.ID,
		When:        jsonHelpers.RFC3339Nano(time.Now().Add(-(time.Minute * 7))),
		ReferenceNr: "b",
	}
}

var (
	// Key1 is a mock api key
	Key1 *models.APIKey
	// Key2 is a mock api key
	Key2 *models.APIKey
	// Key3 is a mock api key
	Key3 *models.APIKey
	// DashboardKey is the mock key for the dashboard
	DashboardKey *models.APIKey
)

var (
	// Profile1 contains the second example profile
	Profile1 *models.Profile
	// Profile2 contains the second example profile
	Profile2 *models.Profile
)

var yearsSinceEducation = 2
var matchedEducation = "MBO 4 ict"
var matchedCourse = "Typecursus"
var professionExperienced = "Dancer"

var (
	// mockMatch1 contains a example match between profile 1 and a cv
	mockMatch1 *models.Match
	mockMatch2 *models.Match
)

// NewMockDB returns an in memory temp testing database with mock data
func NewMockDB() *testingdb.TestConnection {
	conn := testingdb.NewDB()

	// Insert api keys
	conn.UnsafeInsert(
		Key1,
		Key2,
		Key3,
		DashboardKey,
	)

	// Insert secrets
	conn.UnsafeInsert(
		models.UnsafeMustCreateSecret(
			Key1.ID,
			"foo",
			"very-secret-key-of-more-than-16-chars",
			[]byte(`{"foo": 1}`),
			"test secret 1 from mock data",
			models.SecretValueStructureFree,
		),
		models.UnsafeMustCreateSecret(
			Key2.ID,
			"bar",
			"very-secret-key-of-more-than-16-chars",
			[]byte(`{"username": "foo", "password": "bar"}`),
			"test secret 2 from mock data",
			models.SecretValueStructureUser,
		),
	)

	// Insert profiles
	conn.UnsafeInsert(
		Profile1,
		Profile2,
	)

	// Insert matches
	conn.UnsafeInsert(
		mockMatch1,
		mockMatch2,
	)

	return conn
}
