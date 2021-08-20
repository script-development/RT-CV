package mock

import (
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/db/testingdb"
	"github.com/script-development/RT-CV/models"
)

var (
	Key1 = &models.ApiKey{
		M:       dbInterfaces.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "abc",
		Roles:   models.ApiKeyRoleAll,
	}
	Key2 = &models.ApiKey{
		M:       dbInterfaces.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "abc",
		Roles:   models.ApiKeyRoleScraper,
	}
	Key3 = &models.ApiKey{
		M:       dbInterfaces.NewM(),
		Enabled: true,
		Domains: []string{"werk.nl"},
		Key:     "def",
		Roles:   models.ApiKeyRoleInformationObtainer,
	}
)

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
		models.UnsafeMustCreateSecret(Key1.ID, "foo", "very-secret-key", []byte(`{"foo": 1}`)),
		models.UnsafeMustCreateSecret(Key2.ID, "bar", "very-secret-key", []byte(`{"bar": 2}`)),
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
			Zipcodes: []models.Zipcode{{
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
