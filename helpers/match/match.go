package match

import (
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	fuzzymatcher "github.com/mjarkk/fuzzy-matcher"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/helpers/wordvalidator"
	"github.com/script-development/RT-CV/models"
)

// FoundMatch contains a match and why something is matched
type FoundMatch struct {
	Matches models.Match   `json:"matches"`
	Profile models.Profile `json:"profile"`
}

// Match tries to match a profile to a CV
func Match(domains []string, profiles []*models.Profile, cv models.CV) []FoundMatch {
	res := []FoundMatch{}

	formattedDomains := make([][]string, len(domains))
	for idx, domain := range domains {
		formattedDomains[idx] = strings.Split(domain, ".")
	}

	now := time.Now()

	var normalizedCVDriversLicenseCache []string

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
			When:      jsonHelpers.RFC3339Nano(time.Now()),
		}

		// Check domain
		if len(profile.Domains) > 0 && len(domains) > 0 {
			foundMatch := false
			if profile.DomainPartsCache == nil {
				profile.DomainPartsCache = make([][]string, len(profile.Domains))
				for idx, domain := range profile.Domains {
					profile.DomainPartsCache[idx] = strings.Split(domain, ".")
				}
			}
			for domainIdx, domain := range profile.Domains {
				if domain == "*" {
					// This is a match all domain name
					// We can just match the first formatted domain name
					match.Domain = &domains[0]
					foundMatch = true
				}

				domainParts := profile.DomainPartsCache[domainIdx]
				domainPartsLen := len(domainParts)

				for formattedDomainsIdx, formattedDomain := range formattedDomains {
					if len(formattedDomain) == 1 && formattedDomain[0] == "*" {
						// This is match all domain name *
						match.Domain = &domains[formattedDomainsIdx]
						foundMatch = true
						break
					}

					if len(formattedDomain) != domainPartsLen {
						continue
					}

					matched := true
					for i := 0; i < domainPartsLen; i++ {
						domainPart := domainParts[i]
						formattedDomainPart := formattedDomain[i]

						if formattedDomainPart == "*" || domainPart == "*" {
							continue
						}
						if domainPart != formattedDomainPart {
							matched = false
							break
						}
					}

					if matched {
						match.Domain = &domains[formattedDomainsIdx]
						foundMatch = true
						break
					}
				}
			}
			if !foundMatch {
				continue
			}
		}

		// Check years since education
		if profile.YearsSinceEducation > 0 {
			foundMatch := false

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

			if !foundMatch {
				for _, cvCourse := range cv.Courses {
					if len(cvCourse.Name) == 0 || cvCourse.EndDate == nil {
						continue
					}

					t := cvCourse.EndDate.Time()
					if t.After(mustBeAfter) {
						lastEducativeYear = t
					}
				}
			}

			if lastEducativeYear.Before(mustBeAfter) {
				continue
			}

			yearsSinceEducation := time.Now().Year() - lastEducativeYear.Year()
			match.YearsSinceEducation = &yearsSinceEducation
		}

		// Check education and courses
		matchedAnEducationOrCourse := false
		checkedForEducationOrCourse := len(profile.Educations) > 0
		if checkedForEducationOrCourse {
			if (len(cv.Educations) > 0 || len(cv.Courses) > 0) && profile.EducationFuzzyMatcher == nil {
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

			if len(cv.Courses) > 0 {
				for _, cvCourse := range cv.Courses {
					if len(cvCourse.Name) == 0 {
						continue
					}

					educationIdx := profile.EducationFuzzyMatcher.Match(cvCourse.Name)
					if educationIdx == -1 {
						continue
					}

					match.Course = &profile.Educations[educationIdx].Name
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
				profile.NormalizedDriversLicensesCache = []string{}
				for _, l := range profile.DriversLicenses {
					normalizedDriversLicense := wordvalidator.NormalizeString(l.Name)
					if len(normalizedDriversLicense) == 0 {
						continue
					}
					profile.NormalizedDriversLicensesCache = append(profile.NormalizedDriversLicensesCache, normalizedDriversLicense)
				}
			}

		driversLicensesLoop:
			for _, normalizedDriversLicense := range profile.NormalizedDriversLicensesCache {
				// Cache the normalized names once we need it so we don't have to do duplicated work
				if normalizedCVDriversLicenseCache == nil {
					normalizedCVDriversLicenseCache = []string{}
					for _, cvDriversLicense := range cv.DriversLicenses {
						normalizedName := wordvalidator.NormalizeString(cvDriversLicense)
						if len(normalizedName) == 0 {
							continue
						}
						normalizedCVDriversLicenseCache = append(normalizedCVDriversLicenseCache, normalizedName)
					}
				}

				for _, normalizedCVDriversLicense := range normalizedCVDriversLicenseCache {
					if normalizedDriversLicense == normalizedCVDriversLicense {
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
func (match FoundMatch) HandleMatch(cv models.CV, pdfBytes []byte) {
	onMatch := match.Profile.OnMatch

	for _, http := range onMatch.HTTPCall {
		go func(http models.ProfileHTTPCallData) {
			http.MakeRequest(match.Profile, match.Matches)
		}(http)
	}

	emailsLen := len(onMatch.SendMail)
	if emailsLen != 0 {
		emailBody, err := cv.GetEmailHTML(match.Profile, match.Matches.GetMatchSentence())
		if err != nil {
			log.WithError(err).Error("unable to generate email body from CV")
			return
		}

		for _, email := range onMatch.SendMail {
			err := email.SendEmail(match.Profile, emailBody.Bytes(), pdfBytes)
			if err != nil {
				log.WithError(err).Error("unable to send email")
			}
		}
	}
}
