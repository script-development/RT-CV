package mock

import (
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/testingdb"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// Key1 is a mock api key
	Key1 = &models.APIKey{
		M:       db.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "aaa",
		Roles:   models.APIKeyRoleAll,
	}
	// Key2 is a mock api key
	Key2 = &models.APIKey{
		M:       db.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "bbb",
		Roles:   models.APIKeyRoleScraper,
	}
	// Key3 is a mock api key
	Key3 = &models.APIKey{
		M:       db.NewM(),
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
		System:  true,
		Enabled: true,
		Domains: []string{"*"},
		Key:     "ddd",
		Roles:   models.APIKeyRoleDashboard,
	}
)

var (
	// Profile1 contains the first example profile
	Profile1 = &models.Profile{
		Name:                  "Mock profile 1",
		YearsSinceWork:        nil,
		Active:                true,
		MustExpProfession:     true,
		MustDesiredProfession: false,
		MustEducation:         true,
		MustEducationFinished: true,
		MustDriversLicense:    true,
		Domains:               []string{"werk.nl"},
		ListProfile:           true,
		YearsSinceEducation:   1,
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
			SendMail: []models.ProfileSendEmailData{{
				Email: "abc@example.com",
			}},
		},
		Zipcodes: []models.ProfileDutchZipcode{{
			From: 2000,
			To:   8000,
		}},
	}
	// Profile2 contains the second example profile
	Profile2 = &models.Profile{
		Name:                  "Mock profile 2",
		YearsSinceWork:        nil,
		Active:                true,
		MustExpProfession:     false,
		MustDesiredProfession: false,
		MustEducation:         false,
		MustEducationFinished: false,
		MustDriversLicense:    false,
		Domains:               []string{"werk.nl"},
		ListProfile:           false,
		YearsSinceEducation:   0,
		DesiredProfessions:    nil,
		ProfessionExperienced: nil,
		DriversLicenses:       nil,
		Educations:            nil,
		Zipcodes:              nil,
	}
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
		),
		models.UnsafeMustCreateSecret(
			Key2.ID,
			"bar",
			"very-secret-key-of-more-than-16-chars",
			[]byte(`{"bar": 2}`),
			"test secret 2 from mock data",
		),
	)

	// Insert profiles
	conn.UnsafeInsert(
		Profile1,
		Profile2,
	)

	return conn
}
