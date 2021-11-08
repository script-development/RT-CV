package match

import (
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
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
func Match(domains []string, profiles []models.Profile, cv models.CV) []FoundMatch {
	res := []FoundMatch{}

	normalizeString := func(in string) string {
		in = strings.ReplaceAll(in, "\t", " ")
		in = strings.ReplaceAll(in, "\n", " ")

		// Replace double all spaces with normal spaces
		// note that when you have a string with a non even amount of spaces we need to run the replace twice
		in = strings.ReplaceAll(in, "  ", " ")
		in = strings.ReplaceAll(in, "  ", " ")

		return strings.ToLower(strings.TrimSpace(in))
	}

	formattedDomains := make([][]string, len(domains))
	for idx, domain := range domains {
		formattedDomains[idx] = strings.Split(normalizeString(domain), ".")
	}

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

		match := FoundMatch{
			Profile: profile,
			Matches: models.Match{
				M:         db.NewM(),
				ProfileID: profile.ID,
				When:      jsonHelpers.RFC3339Nano(time.Now()),
			},
		}

		// Check domain
		if len(profile.Domains) > 0 && len(domains) > 0 {
			foundMatch := false
			for _, domain := range profile.Domains {
				if domain == "*" {
					// This is a match all domain name
					// We can just match the first formatted domain name
					match.Matches.Domain = &domains[0]
					foundMatch = true
				}

				domainParts := strings.Split(normalizeString(domain), ".")
				domainPartsLen := len(domainParts)

				for formattedDomainsIdx, formattedDomain := range formattedDomains {
					if len(formattedDomain) == 1 && formattedDomain[0] == "*" {
						// This is match all domain name *
						match.Matches.Domain = &domains[formattedDomainsIdx]
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
						match.Matches.Domain = &domains[formattedDomainsIdx]
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
			match.Matches.YearsSinceEducation = &yearsSinceEducation
		}

		// Check education and courses
		matchedAnEducationOrCourse := false
		checkedForEducationOrCourse := len(profile.Educations) > 0
		if checkedForEducationOrCourse {
			if len(cv.Educations) > 0 {
			educationLoop:
				for profileEducationIdx, profileEducation := range profile.Educations {
					if len(profileEducation.Name) == 0 {
						// We don't want those yee yee ass fake educations!
						continue
					}

					for _, cvEducation := range cv.Educations {
						if len(cvEducation.Name) == 0 {
							continue
						}

						if !cvEducation.HasDiploma && profile.MustEducationFinished {
							continue
						}

						if !wordvalidator.IsSame(cvEducation.Name, profileEducation.Name) {
							// Not a equal education title
							continue
						}

						match.Matches.Education = &profile.Educations[profileEducationIdx].Name
						matchedAnEducationOrCourse = true
						break educationLoop
					}
				}
			}

			if len(cv.Courses) > 0 {
			coursesLoop:
				for profileCourseIdx, profileCourse := range profile.Educations {
					if len(profileCourse.Name) == 0 {
						continue
					}

					for _, cvCourse := range cv.Courses {
						if len(cvCourse.Name) == 0 {
							continue
						}

						if !wordvalidator.IsSame(cvCourse.Name, profileCourse.Name) {
							// Not a equal education/course title
							continue
						}

						match.Matches.Course = &profile.Educations[profileCourseIdx].Name
						matchedAnEducationOrCourse = true
						break coursesLoop
					}
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
				profileName := normalizeString(profileProfession.Name)
				if len(profileName) == 0 {
					continue
				}

				for _, cvPreferredJob := range cv.PreferredJobs {
					cvName := normalizeString(cvPreferredJob)
					if len(cvName) == 0 {
						continue
					}

					if cvName == profileName {
						match.Matches.DesiredProfession = &profile.DesiredProfessions[profileProfessionIdx].Name
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
		checkedForProfessionExperienced := len(profile.ProfessionExperienced) > 0
		if checkedForProfessionExperienced {
		professionExperiencedProfileLoop:
			for _, profileProfession := range profile.ProfessionExperienced {
				profileName := normalizeString(profileProfession.Name)
				if len(profileName) == 0 {
					continue
				}

				for _, cvWorkExp := range cv.WorkExperiences {
					profName := normalizeString(cvWorkExp.Profession)
					if len(profName) == 0 {
						continue
					}

					if profName == profileName {
						matchedAProfessionExperienced = true
						break professionExperiencedProfileLoop
					}
				}
			}

			if matchedAProfessionExperienced {
				match.Matches.ProfessionExperienced = true
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
			match.Matches.YearsSinceWork = &yearsSinceLastWork
		}

		// Check drivers license
		matchedADriversLicense := false
		checkedForDriversLicense := len(profile.DriversLicenses) > 0
		if checkedForDriversLicense {
		driversLicensesLoop:
			for _, profileDriversLicense := range profile.DriversLicenses {
				profileName := normalizeString(profileDriversLicense.Name)
				if len(profileName) == 0 {
					continue
				}

				for _, cvDriversLicense := range cv.DriversLicenses {
					cvName := normalizeString(cvDriversLicense)
					if len(cvName) == 0 {
						continue
					}

					if profileName == cvName {
						matchedADriversLicense = true
						break driversLicensesLoop
					}
				}
			}

			if matchedADriversLicense {
				match.Matches.DriversLicense = true
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
					match.Matches.ZipCode = &profile.Zipcodes[idx]
					cvZipInRange = true
					break
				}
			}

			if !cvZipInRange {
				// no matching zipcode
				continue
			}
		}

		res = append(res, match)
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
