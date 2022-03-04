package match

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	fuzzymatcher "github.com/mjarkk/fuzzy-matcher"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/models"
)

// FoundMatch contains a match and why something is matched
type FoundMatch struct {
	Matches models.Match   `json:"matches"`
	Profile models.Profile `json:"profile"`
}

// Match tries to match a profile to a CV
func Match(scraperKey *models.APIKey, profiles []*models.Profile, cv models.CV) []FoundMatch {
	res := []FoundMatch{}

	now := time.Now()

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
			M:         db.NewM(),
			ProfileID: profile.ID,
			When:      jsonHelpers.RFC3339Nano(now),
		}

		// Check domain
		if len(profile.AllowedScrapers) > 0 {
			foundMatch := false
			for _, id := range profile.AllowedScrapers {
				if id == scraperKey.ID {
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				continue
			}
		}

		// Check years since education
		if profile.YearsSinceEducation > 0 {
			mustBeAfter := now.AddDate(-profile.YearsSinceEducation, 0, 0)
			lastEducativeYear := time.Date(1980, time.January, 1, 0, 0, 0, 0, time.Local)

			for _, cvEducation := range cv.Educations {
				if len(cvEducation.Name) == 0 || cvEducation.EndDate == nil {
					continue
				}

				t := cvEducation.EndDate.Time()
				if t.After(lastEducativeYear) {
					lastEducativeYear = t
				}
			}

			if lastEducativeYear.Before(mustBeAfter) {
				continue
			}

			yearsSinceEducation := now.Year() - lastEducativeYear.Year()
			match.YearsSinceEducation = &yearsSinceEducation
		}

		// Check education and courses
		matchedAnEducationOrCourse := false
		checkedForEducationOrCourse := len(profile.Educations) > 0
		if checkedForEducationOrCourse {
			if len(cv.Educations) > 0 && profile.EducationFuzzyMatcher == nil {
				// The fuzzy matcher is not yet setup, lets set it up here
				names := make([]string, len(profile.Educations))
				for idx, education := range profile.Educations {
					names[idx] = education.Name
				}
				profile.EducationFuzzyMatcher = fuzzymatcher.NewMatcher(names...)
			}

			if len(cv.Educations) > 0 {
				for _, cvEducation := range cv.Educations {
					if len(cvEducation.Name) == 0 {
						continue
					}

					if !cvEducation.HasDiploma && profile.MustEducationFinished {
						continue
					}

					educationIdx := profile.EducationFuzzyMatcher.Match(cvEducation.Name)
					if educationIdx == -1 {
						continue
					}

					match.Education = &profile.Educations[educationIdx].Name
					matchedAnEducationOrCourse = true
					break
				}
			}

			if !matchedAnEducationOrCourse && profile.MustEducation {
				// CV doesn't have any matched education
				continue
			}
		}

		// Check profession
		matchedADesiredProfession := false
		checkedForDesiredProfession := len(profile.DesiredProfessions) > 0
		if checkedForDesiredProfession {
			if profile.DesiredProfessionsFuzzyMatcher == nil {
				profileProfessionNames := make([]string, len(profile.DesiredProfessions))
				for idx, p := range profile.DesiredProfessions {
					profileProfessionNames[idx] = p.Name
				}
				profile.DesiredProfessionsFuzzyMatcher = fuzzymatcher.NewMatcher(profileProfessionNames...)
			}

			for _, cvPreferredJob := range cv.PreferredJobs {
				if len(cvPreferredJob) == 0 {
					continue
				}

				matchedDesiredProfession := profile.DesiredProfessionsFuzzyMatcher.Match(cvPreferredJob)
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

		// check profession experienced
		matchedAProfile := false
		checkedForProfessionExperienced := len(profile.ProfessionExperienced) > 0
		if checkedForProfessionExperienced {
			matchedProfileIdx := -1
			for _, workExp := range cv.WorkExperiences {
				if profile.ProfessionExperiencedFuzzyMatcher == nil {
					// The fuzzy matcher is not yet setup, lets set it up here
					names := make([]string, len(profile.ProfessionExperienced))
					for idx, profile := range profile.ProfessionExperienced {
						names[idx] = profile.Name
					}

					profile.ProfessionExperiencedFuzzyMatcher = fuzzymatcher.NewMatcher(names...)
				}

				if len(workExp.Profession) == 0 {
					continue
				}

				match := profile.ProfessionExperiencedFuzzyMatcher.Match(workExp.Profession)
				if match != -1 {
					matchedProfileIdx = match
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
		if profile.YearsSinceWork != nil && *profile.YearsSinceWork > 0 {
			profileMustYearsSinceWork := *profile.YearsSinceWork
			nowYear := now.Year()
			lastWorkYear := 0

			for _, cvWorkExp := range cv.WorkExperiences {
				if cvWorkExp.EndDate == nil {
					continue
				}

				endDateYear := cvWorkExp.EndDate.Time().Year()
				if endDateYear > lastWorkYear {
					lastWorkYear = endDateYear
				}
			}

			// Sanity check
			if lastWorkYear > nowYear {
				lastWorkYear = nowYear
			}

			yearsSinceLastWork := nowYear - lastWorkYear
			if yearsSinceLastWork > profileMustYearsSinceWork {
				// To long ago since last work
				continue
			}

			match.YearsSinceWork = &yearsSinceLastWork
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
			zipStr := strings.TrimSpace(cv.PersonalDetails.Zip)
			zipStrLen := len(zipStr)
			if zipStrLen != 4 && zipStrLen != 6 {
				// Client has invalid zipcode
				continue
			}
			cvZipNr, err := strconv.Atoi(zipStr[:4])
			if err != nil {
				// Client has invalid zipcode
				continue
			}
			cvZipNrUint16 := uint16(cvZipNr)

			cvZipInRange := false
			for idx, zipcode := range profile.Zipcodes {
				if zipcode.IsWithinCithAndArea(cvZipNrUint16) {
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

// HandleMatch sends a match to the desired destination based on the OnMatch field in the profile
func (match FoundMatch) HandleMatch(cv models.CV, pdfFile *os.File, keyName string) {
	onMatch := match.Profile.OnMatch

	for _, http := range onMatch.HTTPCall {
		go func(http models.ProfileHTTPCallData) {
			http.MakeRequest(match.Profile, match.Matches)
		}(http)
	}

	emailsLen := len(onMatch.SendMail)
	if emailsLen != 0 {
		emailBody, err := cv.GetEmailHTML(match.Profile, match.Matches.GetMatchSentence(), keyName)
		if err != nil {
			log.WithError(err).Error("unable to generate email body from CV")
		} else {
			for _, email := range onMatch.SendMail {
				err := email.SendEmail(match.Profile, emailBody.Bytes(), pdfFile)
				if err != nil {
					log.WithError(err).Error("unable to send email")
				}
			}
		}
	}
}
