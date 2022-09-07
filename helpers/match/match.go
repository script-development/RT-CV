package match

import (
	"math"
	"sort"
	"strings"
	"time"

	fuzzymatcher "github.com/mjarkk/fuzzy-matcher"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FoundMatch contains a match and why something is matched
type FoundMatch struct {
	Matches models.Match   `json:"matches"`
	Profile models.Profile `json:"profile"`
}

// Match tries to match a profile to a CV
func Match(scraperKeyID, requestID primitive.ObjectID, profiles []*models.Profile, cv models.CV) []FoundMatch {
	res := []FoundMatch{}

	now := time.Now()
	nowAsMonths := totalMonths(now)

	sort.Sort(cv.Educations)
	sort.Sort(cv.WorkExperiences)

	for _, profile := range profiles {
		if !profile.Active {
			continue
		}

		// There are a lot of CVs that fail on this check on the end
		// Lets make those cases quick as we can easily check that
		if len(profile.Zipcodes) != 0 && len(cv.PersonalDetails.Zip) == 0 {
			continue
		}

		match := models.Match{
			RequestID:   requestID,
			ProfileID:   profile.ID,
			KeyID:       scraperKeyID,
			When:        jsonHelpers.RFC3339Nano(now),
			ReferenceNr: cv.ReferenceNumber,
		}

		// Check domain
		if len(profile.AllowedScrapers) > 0 {
			foundMatch := false
			for _, id := range profile.AllowedScrapers {
				if id == scraperKeyID {
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				continue
			}
		}

		checkYearsSinceEducation := profile.YearsSinceEducation != nil && *profile.YearsSinceEducation > 0

		// Check education and courses
		matchedAnEducationOrCourse := false
		checkedForEducationOrCourse := len(profile.Educations) > 0
		var matchedAnEducationOrCourseDate *time.Time
		if checkedForEducationOrCourse {
			if len(cv.Educations) > 0 && profile.EducationFuzzyMatcherCache == nil {
				// The fuzzy matcher is not yet setup, lets set it up here
				names := make([]string, len(profile.Educations))
				for idx, education := range profile.Educations {
					names[idx] = education.Name
				}
				profile.EducationFuzzyMatcherCache = fuzzymatcher.NewMatcher(names...)
			}

			if len(cv.Educations) > 0 {
				for _, cvEducation := range cv.Educations {
					if len(cvEducation.Name) == 0 {
						continue
					}

					if !cvEducation.HasDiploma && profile.MustEducationFinished {
						continue
					}

					educationIdx := profile.EducationFuzzyMatcherCache.Match(cvEducation.Name)
					if educationIdx == -1 {
						continue
					}

					if !matchedAnEducationOrCourse {
						match.Education = &profile.Educations[educationIdx].Name
						matchedAnEducationOrCourse = true
					}
					if checkYearsSinceEducation {
						date := cvEducation.Date()
						if date != nil {
							dateTime := date.Time()
							matchedAnEducationOrCourseDate = &dateTime
						} else {
							continue
						}
					}
					break
				}
			}

			if !matchedAnEducationOrCourse && profile.MustEducation {
				// CV doesn't have any matched education
				continue
			}
		}

		// Check years since education
		if checkYearsSinceEducation {
			lastEducation := time.Date(1980, time.January, 1, 0, 0, 0, 0, time.Local)

			if matchedAnEducationOrCourseDate != nil {
				lastEducation = *matchedAnEducationOrCourseDate
			} else {
				for _, cvEducation := range cv.Educations {
					date := cvEducation.Date()
					if len(cvEducation.Name) == 0 || date == nil {
						continue
					}

					dateTime := date.Time()
					if dateTime.After(lastEducation) {
						lastEducation = dateTime
					}
				}
			}

			yearsSinceEducation := yearSince(nowAsMonths, totalMonths(lastEducation))
			if yearsSinceEducation > *profile.YearsSinceEducation {
				continue
			}

			match.YearsSinceEducation = &yearsSinceEducation
		}

		// Check profession
		matchedADesiredProfession := false
		checkedForDesiredProfession := len(profile.DesiredProfessions) > 0
		if checkedForDesiredProfession {
			if profile.DesiredProfessionsFuzzyMatcherCache == nil {
				profileProfessionNames := make([]string, len(profile.DesiredProfessions))
				for idx, p := range profile.DesiredProfessions {
					profileProfessionNames[idx] = p.Name
				}
				profile.DesiredProfessionsFuzzyMatcherCache = fuzzymatcher.NewMatcher(profileProfessionNames...)
			}

			for _, cvPreferredJob := range cv.PreferredJobs {
				if len(cvPreferredJob) == 0 {
					continue
				}

				matchedDesiredProfession := profile.DesiredProfessionsFuzzyMatcherCache.Match(cvPreferredJob)
				if matchedDesiredProfession != -1 {
					match.DesiredProfession = &profile.DesiredProfessions[matchedDesiredProfession].Name
					matchedADesiredProfession = true
					break
				}
			}

			if !matchedADesiredProfession && profile.MustDesiredProfession {
				// CV doesn't have any matching professions
				continue
			}
		}

		checkYearsSinceWork := profile.YearsSinceWork != nil && *profile.YearsSinceWork > 0

		// check profession experienced
		matchedAProfile := false
		checkedForProfessionExperienced := len(profile.ProfessionExperienced) > 0
		matchedProfileIdx := -1
		var matchedProfileLastWorkExp *time.Time
		if checkedForProfessionExperienced {
			for _, workExp := range cv.WorkExperiences {
				if profile.ProfessionExperiencedFuzzyMatcherCache == nil {
					// The fuzzy matcher is not yet setup, lets set it up here
					names := make([]string, len(profile.ProfessionExperienced))
					for idx, profile := range profile.ProfessionExperienced {
						names[idx] = profile.Name
					}

					profile.ProfessionExperiencedFuzzyMatcherCache = fuzzymatcher.NewMatcher(names...)
				}

				if len(workExp.Profession) == 0 {
					continue
				}

				match := profile.ProfessionExperiencedFuzzyMatcherCache.Match(workExp.Profession)
				if match == -1 {
					continue
				}

				if matchedProfileIdx == -1 {
					matchedProfileIdx = match
				}

				if !checkYearsSinceWork {
					break
				}

				workExpDate := workExp.Date()
				if workExpDate != nil {
					newValue := workExpDate.Time()
					matchedProfileLastWorkExp = &newValue
					break
				}
			}

			if matchedProfileIdx != -1 {
				matchedAProfile = true
				match.ProfessionExperienced = &profile.ProfessionExperienced[matchedProfileIdx].Name
			} else if profile.MustExpProfession {
				continue
			}
		}

		// Check years since work
		if checkYearsSinceWork {
			profileMustYearsSinceWork := *profile.YearsSinceWork
			lastWorkExp := time.Date(1980, time.January, 1, 0, 0, 0, 0, time.Local)

			if matchedProfileLastWorkExp != nil {
				lastWorkExp = *matchedProfileLastWorkExp
			} else {
				for _, cvWorkExp := range cv.WorkExperiences {
					date := cvWorkExp.Date()
					if date == nil {
						continue
					}
					dateTime := date.Time()
					if dateTime.After(lastWorkExp) {
						lastWorkExp = dateTime
					}
				}
			}

			// Sanity check
			if lastWorkExp.After(now) {
				lastWorkExp = now
			}

			yearsSinceLastWorkExp := yearSince(nowAsMonths, totalMonths(lastWorkExp))
			if yearsSinceLastWorkExp > profileMustYearsSinceWork {
				// To long ago since last work
				continue
			}

			match.YearsSinceWork = &yearsSinceLastWorkExp
		}

		// Check drivers license
		matchedADriversLicense := false
		checkedForDriversLicense := len(profile.DriversLicenses) > 0
		if checkedForDriversLicense {
			if profile.NormalizedDriversLicensesCache == nil {
				profile.NormalizedDriversLicensesCache = []jsonHelpers.DriversLicense{}
				for _, l := range profile.DriversLicenses {
					normalizedDriversLicense := strings.ToUpper(strings.ReplaceAll(l.Name, " ", ""))
					if len(normalizedDriversLicense) == 0 {
						continue
					}
					profile.NormalizedDriversLicensesCache = append(
						profile.NormalizedDriversLicensesCache,
						jsonHelpers.NewDriversLicense(normalizedDriversLicense),
					)
				}
			}

		driversLicensesLoop:
			for _, normalizedDriversLicense := range profile.NormalizedDriversLicensesCache {
				for _, cvDriversLicense := range cv.DriversLicenses {
					if normalizedDriversLicense == cvDriversLicense {
						matchedADriversLicense = true
						break driversLicensesLoop
					}
				}
			}

			if matchedADriversLicense {
				match.DriversLicense = true
			} else if profile.MustDriversLicense {
				// CV doesn't have any matching drivers license
				continue
			}
		}

		// Check if at least one of the matches is true
		if checkedForEducationOrCourse || checkedForDesiredProfession || checkedForDriversLicense || checkedForProfessionExperienced {
			if !matchedAnEducationOrCourse && !matchedADesiredProfession && !matchedADriversLicense && !matchedAProfile {
				continue
			}
		}

		// Check zipcodes
		if len(profile.Zipcodes) != 0 {
			cvZip, validCVZip := cv.PersonalDetails.ZipAsNr()
			if !validCVZip {
				continue
			}

			cvZipInRange := false
			for idx, zipcode := range profile.Zipcodes {
				if zipcode.IsWithinCithAndArea(cvZip) {
					match.ZipCode = &profile.Zipcodes[idx]
					cvZipInRange = true
					break
				}
			}

			if !cvZipInRange {
				// no matching zipcode
				continue
			}
		}

		res = append(res, FoundMatch{
			Profile: *profile,
			Matches: match,
		})
	}

	return res
}

func totalMonths(t time.Time) int {
	return t.Year()*12 + int(t.Month()) - 1
}

func yearSince(nowInMonths int, comparedToMonths int) int {
	return int(math.Round(float64(nowInMonths-comparedToMonths) / 12))
}
