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

	var (
		// normalizedCVEducationCache      []string
		normalizedCVPreferredJobsCache  []string
		normalizedCVWrkExpProfnameCache []string
		normalizedCVDriversLicenseCache []string
	)

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
			for _, domain := range profile.Domains {
				if domain == "*" {
					// This is a match all domain name
					// We can just match the first formatted domain name
					match.Domain = &domains[0]
					foundMatch = true
				}

				domainParts := strings.Split(domain, ".")
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
				educationNames := make([]string, len(profile.Educations))
				for idx, education := range profile.Educations {
					educationNames[idx] = education.Name
				}
				profile.EducationFuzzyMatcher = fuzzymatcher.NewMatcher(educationNames...)
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
		professionLoop:
			for profileProfessionIdx, profileProfession := range profile.DesiredProfessions {
				profileName := wordvalidator.NormalizeString(profileProfession.Name)
				if len(profileName) == 0 {
					continue
				}

				// Cache the normalized names once we need it so we don't have to do duplicated work
				if normalizedCVPreferredJobsCache == nil {
					normalizedCVPreferredJobsCache = []string{}
					for _, cvPreferredJob := range cv.PreferredJobs {
						normalizedName := wordvalidator.NormalizeString(cvPreferredJob)
						if len(normalizedName) == 0 {
							continue
						}
						normalizedCVPreferredJobsCache = append(normalizedCVPreferredJobsCache, normalizedName)
					}
				}

				for _, cvName := range normalizedCVPreferredJobsCache {
					if cvName == profileName {
						match.DesiredProfession = &profile.DesiredProfessions[profileProfessionIdx].Name
						matchedADesiredProfession = true
						break professionLoop
					}
				}
			}

			if !matchedADesiredProfession && profile.MustDesiredProfession {
				// CV doesn't have any matching professions
				continue
			}
		}

		// check profession experienced
		matchedAProfessionExperienced := false
		matchedProfileName := ""
		checkedForProfessionExperienced := len(profile.ProfessionExperienced) > 0
		if checkedForProfessionExperienced {
		professionExperiencedProfileLoop:
			for _, profileProfession := range profile.ProfessionExperienced {
				profileName := wordvalidator.NormalizeString(profileProfession.Name)
				if len(profileName) == 0 {
					continue
				}

				// Cache the normalized names once we need it so we don't have to do duplicated work
				if normalizedCVWrkExpProfnameCache == nil {
					normalizedCVWrkExpProfnameCache = []string{}
					for _, cvWorkExp := range cv.WorkExperiences {
						normalizedName := wordvalidator.NormalizeString(cvWorkExp.Profession)
						if len(normalizedName) == 0 {
							continue
						}
						normalizedCVWrkExpProfnameCache = append(normalizedCVWrkExpProfnameCache, normalizedName)
					}
				}

				for _, profName := range normalizedCVWrkExpProfnameCache {
					if profName == profileName {
						matchedAProfessionExperienced = true
						matchedProfileName = profileProfession.Name
						break professionExperiencedProfileLoop
					}
				}
			}

			if matchedAProfessionExperienced {
				match.ProfessionExperienced = &matchedProfileName
			} else if profile.MustExpProfession {
				continue
			}
		}

		// Check years since work
		yearsSinceWork := profile.YearsSinceWork
		if yearsSinceWork != nil && *yearsSinceWork > 0 {
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

			if now.Year()-*yearsSinceWork > lastWorkYear {
				// To long ago since last work
				continue
			}

			yearsSinceLastWork := now.Year() - lastWorkYear
			match.YearsSinceWork = &yearsSinceLastWork
		}

		// Check drivers license
		matchedADriversLicense := false
		checkedForDriversLicense := len(profile.DriversLicenses) > 0
		if checkedForDriversLicense {
		driversLicensesLoop:
			for _, profileDriversLicense := range profile.DriversLicenses {
				profileName := wordvalidator.NormalizeString(profileDriversLicense.Name)
				if len(profileName) == 0 {
					continue
				}

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

				for _, cvName := range normalizedCVDriversLicenseCache {
					if profileName == cvName {
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

		if checkedForEducationOrCourse || checkedForDesiredProfession || checkedForDriversLicense || checkedForProfessionExperienced {
			// Check if at least one of the matches is true
			if !matchedAnEducationOrCourse && !matchedADesiredProfession && !matchedADriversLicense && !matchedAProfessionExperienced {
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
