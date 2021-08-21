package mock

import (
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/db/testingdb"
	"github.com/script-development/RT-CV/models"
)

var (
	// Key1 is a mock api key
	Key1 = &models.APIKey{
		M:       dbInterfaces.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "abc",
		Roles:   models.APIKeyRoleAll,
	}
	// Key2 is a mock api key
	Key2 = &models.APIKey{
		M:       dbInterfaces.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "abc",
		Roles:   models.APIKeyRoleScraper,
	}
	// Key3 is a mock api key
	Key3 = &models.APIKey{
		M:       dbInterfaces.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "def",
		Roles:   models.APIKeyRoleInformationObtainer,
	}
)

// NewMockDB returns an in memory temp testing database with mock data
func NewMockDB() *testingdb.TestConnection {
	db := testingdb.NewDB()

	// Insert api keys
	db.UnsafeInsert(
		Key1,
		Key2,
		Key3,
	)

	// Insert secrets
	db.UnsafeInsert(
		models.UnsafeMustCreateSecret(Key1.ID, "foo", "very-secret-key-of-more-than-16-chars", []byte(`{"foo": 1}`)),
		models.UnsafeMustCreateSecret(Key2.ID, "bar", "very-secret-key-of-more-than-16-chars", []byte(`{"bar": 2}`)),
	)

	// Insert profiles
	db.UnsafeInsert(
		&models.Profile{
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
			DesiredProfessions: []models.Profession{{
				Name: "Rapper",
			}},
			ProfessionExperienced: []models.Profession{{
				Name: "Dancer",
			}},
			DriversLicenses: []models.DriversLicense{{
				Name: "A",
			}},
			Educations: []models.DBEducation{{
				Name: "Default",
			}},
			Emails: []models.Email{{
				Email: "abc@example.com",
			}},
			Zipcodes: []models.DutchZipcode{{
				From: 2000,
				To:   8000,
			}},
		},
		&models.Profile{
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
			Emails:                nil,
			Zipcodes:              nil,
		},
	)

	return db
}
