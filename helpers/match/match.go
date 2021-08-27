package match

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/script-development/RT-CV/helpers/wordvalidator"
	"github.com/script-development/RT-CV/models"
)

// AMatch contains a match and why something is matched
type AMatch struct {
	Matches Matches        `json:"matches"`
	Profile models.Profile `json:"profile"`
}

// Matches contains the areas of the profile where the match was found
type Matches struct {
	Domain                *string                     `json:"domains"`
	YearsSinceWork        bool                        `json:"yearsSinceWork"`
	YearsSinceEducation   bool                        `json:"yearsSinceEducation"`
	EducationOrCourse     bool                        `json:"educationOrCourse"`
	DesiredProfession     bool                        `json:"desiredProfession"`
	ProfessionExperienced bool                        `json:"professionExperienced"`
	DriversLicense        bool                        `json:"driversLicense"`
	ZipCode               *models.ProfileDutchZipcode `json:"zipCode"`
}

// GetMatchSentence returns a
func (m Matches) GetMatchSentence() string {
	sentences := []string{}
	if m.Domain != nil {
		sentences = append(sentences, "domain naam "+*m.Domain)
	}
	if m.YearsSinceWork {
		sentences = append(sentences, "jaren sinds werk")
	}
	if m.YearsSinceEducation {
		sentences = append(sentences, "jaren sinds laatste opleiding")
	}
	if m.EducationOrCourse {
		sentences = append(sentences, "opleiding of cursus")
	}
	if m.DesiredProfession {
		sentences = append(sentences, "gewenste werkveld")
	}
	if m.ProfessionExperienced {
		sentences = append(sentences, "wil profession")
	}
	if m.DriversLicense {
		sentences = append(sentences, "rijbewijs")
	}
	if m.ZipCode != nil {
		sentences = append(sentences, fmt.Sprintf("postcode in range %d - %d", m.ZipCode.From, m.ZipCode.To))
	}

	switch len(sentences) {
	case 0:
		return ""
	case 1:
		return sentences[0]
	default:
		return fmt.Sprintf("%s en %s", strings.Join(sentences[:len(sentences)-1], ", "), sentences[len(sentences)-1])
	}
}

// Match tries to match a profile to a CV
// FIXME: There are a lot of performance optimizations that could be done here
func Match(domains []string, profiles []models.Profile, cv models.CV) []AMatch {
	res := []AMatch{}

	normalizeString := func(in string) string {
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

		match := AMatch{
			Profile: profile,
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

			for _, cvEducation := range cv.Educations {
				if len(cvEducation.Name) == 0 || len(cvEducation.EndDate) == 0 {
					continue
				}

				t, err := time.Parse(time.RFC3339, cvEducation.EndDate)
				if err != nil {
					continue
				}

				if t.AddDate(profile.YearsSinceEducation, 0, 0).After(now) {
					foundMatch = true
					break
				}
			}

			if !foundMatch {
				for _, cvCourse := range cv.Courses {
					if len(cvCourse.Name) == 0 || len(cvCourse.EndDate) == 0 {
						continue
					}

					t, err := time.Parse(time.RFC3339, cvCourse.EndDate)
					if err != nil {
						continue
					}

					if t.AddDate(profile.YearsSinceEducation, 0, 0).After(now) {
						foundMatch = true
						break
					}
				}
			}

			if !foundMatch {
				continue
			}

			match.Matches.YearsSinceEducation = true
		}

		// Check education and courses
		matchedAnEducationOrCourse := false
		checkedForEducationOrCourse := len(profile.Educations) > 0
		if checkedForEducationOrCourse {
			if len(cv.Educations) > 0 {
			educationLoop:
				for _, profileEducation := range profile.Educations {
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

						matchedAnEducationOrCourse = true
						break educationLoop
					}
				}
			}

			if !matchedAnEducationOrCourse && len(cv.Courses) > 0 {
			coursesLoop:
				for _, profileCourse := range profile.Educations {
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

						matchedAnEducationOrCourse = true
						break coursesLoop
					}
				}
			}

			if matchedAnEducationOrCourse {
				match.Matches.EducationOrCourse = true
			} else if profile.MustEducation {
				// CV doesn't have any matched education
				continue
			}
		}

		// Check profession
		matchedADesiredProfession := false
		checkedForDesiredProfession := len(profile.DesiredProfessions) > 0
		if checkedForDesiredProfession {
		professionLoop:
			for _, profileProfession := range profile.DesiredProfessions {
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
						matchedADesiredProfession = true
						break professionLoop
					}
				}
			}

			if matchedADesiredProfession {
				match.Matches.DesiredProfession = true
			} else if profile.MustDesiredProfession {
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
				if cvWorkExp.EndDate == "" {
					continue
				}

				endDate, err := time.Parse(time.RFC3339, cvWorkExp.EndDate)
				if err != nil {
					continue
				}

				endDateYear := endDate.Year()
				if endDateYear > lastWorkYear {
					lastWorkYear = endDateYear
				}
			}

			if now.Year()-*yearsSinceWork > lastWorkYear {
				// To long ago since last work
				continue
			}

			match.Matches.YearsSinceWork = true
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
			zipStr := cv.PersonalDetails.Zip
			if len(zipStr) < 4 {
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
				checkFrom := zipcode.From
				checkTo := zipcode.To
				if checkFrom == 0 && checkTo == 0 {
					continue
				}

				if checkFrom > checkTo {
					// Swap from and to
					originalFrom := checkFrom
					checkFrom = checkTo
					checkTo = originalFrom
				}

				if cvZipNrUint16 >= checkFrom && cvZipNrUint16 <= checkTo {
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
