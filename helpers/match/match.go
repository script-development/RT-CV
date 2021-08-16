package match

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/script-development/RT-CV/helpers/wordvalidator"
	"github.com/script-development/RT-CV/models"
)

// Match tries to match a profile to a CV
// FIXME: There are a lot of performance optimizations that could be done here
func Match(domain string, profiles []models.Profile, cv models.Cv) []models.Profile {
	res := []models.Profile{}

	normalizeString := func(in string) string {
		return strings.ToLower(strings.TrimSpace(in))
	}

	formattedDomain := normalizeString(domain)
	now := time.Now()

	for _, profile := range profiles {
		// Check domain
		if profile.SiteId != nil && *profile.SiteId > 0 && len(formattedDomain) != 0 {
			if normalizeString(profile.Site.Domain) != formattedDomain {
				// Domain doesn't match
				continue
			}
		}

		if !profile.Active {
			continue
		}

		// Check years since education
		if profile.YearsSinceEducation > 0 {
			match := false

			for _, cvEducation := range cv.Educations {
				if len(cvEducation.Name) == 0 || len(cvEducation.EndDate) == 0 {
					continue
				}

				t, err := time.Parse(time.RFC3339, cvEducation.EndDate)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if t.AddDate(profile.YearsSinceEducation, 0, 0).After(now) {
					match = true
					break
				}
			}

			if !match {
				for _, cvCourse := range cv.Courses {
					if len(cvCourse.Name) == 0 || len(cvCourse.EndDate) == 0 {
						continue
					}

					t, err := time.Parse(time.RFC3339, cvCourse.EndDate)
					if err != nil {
						continue
					}

					if t.AddDate(profile.YearsSinceEducation, 0, 0).After(now) {
						match = true
						break
					}
				}
			}

			if !match {
				continue
			}
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

			if profile.MustDesiredProfession && !matchedADesiredProfession {
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

			if profile.MustExpProfession && !matchedAProfessionExperienced {
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

			if profile.MustDriversLicense && !matchedADriversLicense {
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
			for _, zipcode := range profile.Zipcodes {
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
					cvZipInRange = true
					break
				}
			}

			if !cvZipInRange {
				// no matching zipcode
				continue
			}
		}

		res = append(res, profile)
	}

	return res
}
