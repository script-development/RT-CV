package match

import (
	"testing"
	"time"

	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MustMatchSingle(t *testing.T, p models.Profile, cv models.CV) {
	p.AllowedScrapers = []primitive.ObjectID{mock.Key2.ID}
	p.Active = true

	matches := Match(mock.Key2.ID, []*models.Profile{&p}, cv)
	Equal(t, 1, len(matches), matches)
}

func MustNotMatchSingle(t *testing.T, p models.Profile, cv models.CV) {
	p.AllowedScrapers = []primitive.ObjectID{mock.Key2.ID}
	p.Active = true

	matches := Match(mock.Key2.ID, []*models.Profile{&p}, cv)
	Equal(t, 0, len(matches), matches)
}

func TestMatchSiteMismatch(t *testing.T) {
	matches := Match(mock.Key2.ID, []*models.Profile{{
		AllowedScrapers: []primitive.ObjectID{mock.Key1.ID},
		Active:          true,
	}}, models.CV{})
	Equal(t, 0, len(matches), matches)
}

func TestMatchNonActive(t *testing.T) {
	matches := Match(mock.Key2.ID, []*models.Profile{{Active: false}}, models.CV{})
	Equal(t, 0, len(matches), matches)
}

func TestMatchEmptyProfile(t *testing.T) {
	MustMatchSingle(t, models.Profile{}, models.CV{})
}

func TestMatchZipCode(t *testing.T) {
	cases := []string{"1500AB", "1000AB", "2000"}
	for _, caseItem := range cases {
		// Valid 1 to 1 match
		MustMatchSingle(
			t,
			models.Profile{Zipcodes: []models.ProfileDutchZipcode{{From: 1000, To: 2000}}},
			models.CV{PersonalDetails: models.PersonalDetails{Zip: caseItem}},
		)

		// Outside of range
		MustNotMatchSingle(
			t,
			models.Profile{Zipcodes: []models.ProfileDutchZipcode{{From: 6000, To: 9000}}},
			models.CV{PersonalDetails: models.PersonalDetails{Zip: caseItem}},
		)
	}

	// invalid CV zip code
	MustNotMatchSingle(
		t,
		models.Profile{Zipcodes: []models.ProfileDutchZipcode{{From: 1000, To: 2000}}},
		models.CV{PersonalDetails: models.PersonalDetails{Zip: "AAAAAA"}},
	)

	// Multiple zip codes
	MustMatchSingle(
		t,
		models.Profile{Zipcodes: []models.ProfileDutchZipcode{
			{From: 1000, To: 2000},
			{From: 3000, To: 3500},
			{From: 4000, To: 5000},
			{From: 6000, To: 8000},
		}},
		models.CV{PersonalDetails: models.PersonalDetails{Zip: "4100AB"}},
	)

	// Reverse zip code
	MustMatchSingle(
		t,
		models.Profile{Zipcodes: []models.ProfileDutchZipcode{{From: 6000, To: 2000}}},
		models.CV{PersonalDetails: models.PersonalDetails{Zip: "4100AB"}},
	)
}

func TestMatchEducation(t *testing.T) {
	// No educations in CV
	MustNotMatchSingle(
		t,
		models.Profile{MustEducation: true, Educations: []models.ProfileEducation{{}}},
		models.CV{},
	)

	// No educations in CV
	MustNotMatchSingle(
		t,
		models.Profile{Educations: []models.ProfileEducation{{}}},
		models.CV{},
	)

	// Match on education
	MustMatchSingle(
		t,
		models.Profile{Educations: []models.ProfileEducation{{Name: "Bananenplukker"}}},
		models.CV{Educations: []models.Education{{Name: "Bananenplukker"}}},
	)

	// Match with multiple educations
	MustMatchSingle(
		t,
		models.Profile{Educations: []models.ProfileEducation{
			{Name: "professioneel peren eten"},
			{Name: "Bananenplukker"},
		}},
		models.CV{Educations: []models.Education{
			{Name: "Pro gangster"},
			{Name: "Bananenplukker"},
		}},
	)
}

func TestMatchEducationMustFinish(t *testing.T) {
	// Education not finished
	MustNotMatchSingle(
		t,
		models.Profile{
			MustEducationFinished: true,
			Educations:            []models.ProfileEducation{{Name: "Bananenplukker"}},
		},
		models.CV{Educations: []models.Education{{Name: "Bananenplukker"}}},
	)

	// Education finished
	MustMatchSingle(
		t,
		models.Profile{
			MustEducationFinished: true,
			Educations:            []models.ProfileEducation{{Name: "Bananenplukker"}},
		},
		models.CV{Educations: []models.Education{{Name: "Bananenplukker", HasDiploma: true}}},
	)
}

func TestMatchEducationYearsSinceEducation(t *testing.T) {
	MustMatchSingle(
		t,
		models.Profile{
			YearsSinceEducation: 2,
		},
		models.CV{
			Educations: []models.Education{{
				Name:    "Bananenplukker",
				EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-1, 0, 0)).ToPtr(),
			}},
		},
	)

	MustMatchSingle(
		t,
		models.Profile{
			YearsSinceEducation: 2,
		},
		models.CV{
			Educations: []models.Education{
				{
					Name:    "Bananenplukker",
					EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-3, 0, 0)).ToPtr(),
				},
				{
					Name:    "Bananenplukker",
					EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-1, 0, 0)).ToPtr(),
				},
				{
					Name:    "Bananenplukker",
					EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-3, 0, 0)).ToPtr(),
				},
			},
		},
	)

	MustNotMatchSingle(
		t,
		models.Profile{
			YearsSinceEducation: 1,
		},
		models.CV{
			Educations: []models.Education{{
				Name:    "Bananenplukker",
				EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-2, 0, 0)).ToPtr(),
			}},
		},
	)
}

func TestMatchDesiredProfession(t *testing.T) {
	// Match on desired profession
	MustMatchSingle(
		t,
		models.Profile{
			MustDesiredProfession: true,
			DesiredProfessions:    []models.ProfileProfession{{Name: "Bananenplukker"}},
		},
		models.CV{PreferredJobs: []string{"Bananenplukker"}},
	)

	// no desired profession match
	MustNotMatchSingle(
		t,
		models.Profile{
			MustDesiredProfession: true,
			DesiredProfessions:    []models.ProfileProfession{{Name: "Real gangster"}},
		},
		models.CV{PreferredJobs: []string{"Bananenplukker"}},
	)
}

func TestMatchDesiredProfessionExperienced(t *testing.T) {
	// Match on profession experienced
	MustMatchSingle(
		t,
		models.Profile{
			MustExpProfession:     true,
			ProfessionExperienced: []models.ProfileProfession{{Name: "Bananenplukker"}},
		},
		models.CV{WorkExperiences: []models.WorkExperience{{Profession: "Bananenplukker"}}},
	)

	// No profession experienced match
	MustNotMatchSingle(
		t,
		models.Profile{
			MustExpProfession:     true,
			ProfessionExperienced: []models.ProfileProfession{{Name: "Real gangster stuff"}},
		},
		models.CV{WorkExperiences: []models.WorkExperience{{Profession: "Bananenplukker"}}},
	)
}

func TestMatchYearsSinceWork(t *testing.T) {
	yearsSinceWork := 2
	MustMatchSingle(
		t,
		models.Profile{YearsSinceWork: &yearsSinceWork},
		models.CV{WorkExperiences: []models.WorkExperience{{EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-1, 0, 0)).ToPtr()}}},
	)

	yearsSinceWork = 1
	MustNotMatchSingle(
		t,
		models.Profile{YearsSinceWork: &yearsSinceWork},
		models.CV{WorkExperiences: []models.WorkExperience{{EndDate: jsonHelpers.RFC3339Nano(time.Now().AddDate(-2, 0, 0)).ToPtr()}}},
	)
}

func TestMatchDriversLicense(t *testing.T) {
	// Match on drivers license
	MustMatchSingle(
		t,
		models.Profile{
			MustDriversLicense: true,
			DriversLicenses:    []models.ProfileDriversLicense{{Name: "A"}},
		},
		models.CV{DriversLicenses: []jsonHelpers.DriversLicense{jsonHelpers.NewDriversLicense("A")}},
	)

	// No drivers license match
	MustNotMatchSingle(
		t,
		models.Profile{
			MustDriversLicense: true,
			DriversLicenses:    []models.ProfileDriversLicense{{Name: "A"}},
		},
		models.CV{DriversLicenses: []jsonHelpers.DriversLicense{jsonHelpers.NewDriversLicense("B")}},
	)
}

func TestGetMatchSentence(t *testing.T) {
	sentence := (&models.Match{}).GetMatchSentence()
	Equal(t, "", sentence)

	yearsSinceWork := 3
	sentence = (&models.Match{YearsSinceWork: &yearsSinceWork}).GetMatchSentence()
	Equal(t, "3 jaren sinds laatste werkervaring", sentence)

	sentence = (&models.Match{YearsSinceWork: &yearsSinceWork, YearsSinceEducation: &yearsSinceWork}).GetMatchSentence()
	Equal(t, "3 jaren sinds laatste werkervaring en 3 jaren sinds laatste opleiding", sentence)

	zipCode := models.ProfileDutchZipcode{
		From: 2000,
		To:   5000,
	}
	education := "beeing smart"
	profession := "gangster"

	sentence = (&models.Match{
		YearsSinceWork:        &yearsSinceWork,
		YearsSinceEducation:   &yearsSinceWork,
		Education:             &education,
		DesiredProfession:     &profession,
		ProfessionExperienced: &profession,
		DriversLicense:        true,
		ZipCode:               &zipCode,
	}).GetMatchSentence()
	expectedResult := "3 jaren sinds laatste werkervaring" +
		", 3 jaren sinds laatste opleiding" +
		", opleiding beeing smart" +
		", gewenste werkveld gangster" +
		", gewerkt als gangster" +
		", gewenste rijbewijs" +
		" en postcode in range 2000 - 5000"
	Equal(t, expectedResult, sentence)
}

func TestTotalMonths(t *testing.T) {
	now := time.Now()
	totalMonths := totalMonths(now)
	Greater(t, totalMonths, now.Year()*12)
}

func TestYearSince(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		comparedTo time.Time
		expect     int
	}{
		{now, 0},
		{now.AddDate(-1, 0, 0), 1},
		{now.AddDate(-2, 0, 0), 2},
		{now.AddDate(-5, 0, 0), 5},
		{now.AddDate(-1, -6, 0), 2},
		{now.AddDate(0, -6, 0), 1},
	}

	for _, testCase := range testCases {
		Equal(
			t,
			testCase.expect,
			yearSince(totalMonths(now), totalMonths(testCase.comparedTo)),
		)
	}
}
