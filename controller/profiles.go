package controller

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	ctxPkg "github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeAllProfiles = routeBuilder.R{
	Description: "get all profiles stored in the database",
	Res:         []models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		profiles, err := models.GetProfiles(ctxPkg.Get(c).DBConn, nil)
		if err != nil {
			return err
		}
		return c.JSON(profiles)
	},
}

const exampleQuery = `{"labels": {"key": "value", "other_key": {"$exists": true}}}`

var routeQueryProfiles = routeBuilder.R{
	Description: strings.Join([]string{
		"make a MongoDB query directly against the database.",
		"For more info about mongodb queries you can take a look at: https://www.mongodb.com/docs/manual/tutorial/query-documents/#std-label-read-operations-query-argument",
		"An example query would be something like: `" + exampleQuery + "`",
		"Note that you can't filter for the _id field",
	}, "\n\n"),
	Body: primitive.M{},
	Res:  []models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		body := primitive.M{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles, err := models.GetProfiles(ctxPkg.Get(c).DBConn, body)
		if err != nil {
			return err
		}
		return c.JSON(profiles)
	},
}

// RouteGetProfilesCountRes is the response for routeGetProfilesCount
type RouteGetProfilesCountRes struct {
	// Total is the total number of profiles
	Total uint64 `json:"total"`
	// Usable is the number of profiles that can be used for matching
	Usable uint64 `json:"usable"`
}

var routeGetProfilesCount = routeBuilder.R{
	Description: "get the number of profiles stored in the database",
	Res:         RouteGetProfilesCountRes{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctxPkg.Get(c)

		profilesCount, err := models.GetProfilesCount(ctx.DBConn)
		if err != nil {
			return err
		}

		usableProfilesCount, err := models.GetActualMatchActiveProfilesCount(ctx.DBConn)
		if err != nil {
			return err
		}

		res := RouteGetProfilesCountRes{
			Total:  profilesCount,
			Usable: usableProfilesCount,
		}
		return c.JSON(res)
	},
}

func middlewareBindProfile() routeBuilder.M {
	return routeBuilder.M{
		Fn: func(c *fiber.Ctx) error {
			profileParam := c.Params(`profile`)

			profileID, err := primitive.ObjectIDFromHex(profileParam)
			if err != nil {
				return err
			}

			ctx := ctxPkg.Get(c)

			profile, err := models.GetProfile(ctx.DBConn, profileID)
			if err != nil {
				return err
			}

			ctx.Profile = &profile

			return c.Next()
		},
	}
}

var routeGetProfile = routeBuilder.R{
	Description: "get a profile based on it's ID from the database",
	Res:         models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		return c.JSON(ctxPkg.Get(c).Profile)
	},
}

var routeCreateProfile = routeBuilder.R{
	Description: "create a new profile that can match scraped CVs",
	Res:         models.Profile{},
	Body:        models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctxPkg.Get(c)

		var profile models.Profile
		err := c.BodyParser(&profile)
		if err != nil {
			return err
		}

		err = profile.ValidateCreateNewProfile(ctx.DBConn)
		if err != nil {
			return err
		}

		// Set the ID of the profile
		profile.M = db.NewM()

		// Save the profile to the database
		err = ctx.DBConn.Insert(&profile)
		if err != nil {
			return err
		}

		// Invalidate profiles cache
		ctx.ResetMatcherProfilesCache()

		return c.JSON(profile)
	},
}

// UpdateProfileReq are all the fields that can be updated in a profile
// All the top level fields are optional and if null should not be updated
type UpdateProfileReq struct {
	Name            *string  `json:"name"`
	Active          *bool    `json:"active"`
	AllowedScrapers []string `json:"allowedScrapers" description:"the scraper IDs that can be used to match this profile, if null/undefined this value won't be updated, if empty array all scrapers will be allowed"`

	MustDesiredProfession *bool                      `json:"mustDesiredProfession"`
	DesiredProfessions    []models.ProfileProfession `json:"desiredProfessions"`

	YearsSinceWork *int `json:"yearsSinceWork"`

	MustExpProfession     *bool                      `json:"mustExpProfession"`
	ProfessionExperienced []models.ProfileProfession `json:"professionExperienced"`

	MustDriversLicense *bool                          `json:"mustDriversLicense"`
	DriversLicenses    []models.ProfileDriversLicense `json:"driversLicenses"`

	MustEducationFinished *bool                     `json:"mustEducationFinished"`
	MustEducation         *bool                     `json:"mustEducation"`
	YearsSinceEducation   *int                      `json:"yearsSinceEducation"`
	Educations            []models.ProfileEducation `json:"educations"`

	Zipcodes []models.ProfileDutchZipcode `json:"zipCodes"`

	ListsAllowed *bool `json:"listsAllowed"`

	Lables map[string]any `json:"labels" description:"custom labels that can be used by API users to identify profiles, the key needs to be a string and the value can be anything"`

	OnMatch *models.ProfileOnMatch `json:"onMatch"`
}

var routeModifyProfile = routeBuilder.R{
	Description: strings.Join([]string{
		"modify an existing profile.",
		"All the top level body fields are optional thus you only have to provide the fields you want to update.",
	}, "\n\n"),
	Res:  models.Profile{},
	Body: UpdateProfileReq{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctxPkg.Get(c)

		var body UpdateProfileReq
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		if body.Name != nil {
			ctx.Profile.Name = *body.Name
		}
		if body.Active != nil {
			ctx.Profile.Active = *body.Active
		}
		if body.AllowedScrapers != nil {
			allowedScrapersIDs := make([]primitive.ObjectID, len(body.AllowedScrapers))
			for idx, id := range body.AllowedScrapers {
				bsonID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return err
				}
				allowedScrapersIDs[idx] = bsonID
			}

			err = models.CheckAPIKeysExists(ctx.DBConn, allowedScrapersIDs)
			if err != nil {
				return err
			}
			ctx.Profile.AllowedScrapers = allowedScrapersIDs
		}
		if body.MustDesiredProfession != nil {
			ctx.Profile.MustDesiredProfession = *body.MustDesiredProfession
		}
		if body.DesiredProfessions != nil {
			ctx.Profile.DesiredProfessions = body.DesiredProfessions
		}

		if body.YearsSinceWork != nil {
			if *body.YearsSinceWork == 0 {
				ctx.Profile.YearsSinceWork = nil
			} else {
				ctx.Profile.YearsSinceWork = &*body.YearsSinceWork
			}
		}
		if body.MustExpProfession != nil {
			ctx.Profile.MustExpProfession = *body.MustExpProfession
		}
		if body.ProfessionExperienced != nil {
			ctx.Profile.ProfessionExperienced = body.ProfessionExperienced
		}
		if body.MustDriversLicense != nil {
			ctx.Profile.MustDriversLicense = *body.MustDriversLicense
		}
		if body.DriversLicenses != nil {
			ctx.Profile.DriversLicenses = body.DriversLicenses
		}
		if body.MustEducationFinished != nil {
			ctx.Profile.MustEducationFinished = *body.MustEducationFinished
		}
		if body.MustEducation != nil {
			ctx.Profile.MustEducation = *body.MustEducation
		}
		if body.YearsSinceEducation != nil {
			if *body.YearsSinceEducation == 0 {
				ctx.Profile.YearsSinceEducation = nil
			} else {
				ctx.Profile.YearsSinceEducation = &*body.YearsSinceEducation
			}
		}
		if body.Educations != nil {
			ctx.Profile.Educations = body.Educations
		}
		if body.Zipcodes != nil {
			ctx.Profile.Zipcodes = body.Zipcodes
		}
		if body.OnMatch != nil {
			ctx.Profile.OnMatch = *body.OnMatch
		}
		if body.Lables != nil {
			ctx.Profile.Lables = body.Lables
		}
		if body.ListsAllowed != nil {
			ctx.Profile.ListsAllowed = *body.ListsAllowed
		}

		err = ctx.DBConn.UpdateByID(ctx.Profile)
		if err != nil {
			return err
		}

		// Invalidate profiles cache
		ctx.ResetMatcherProfilesCache()

		return c.JSON(ctx.Profile)
	},
}

var routeDeleteProfile = routeBuilder.R{
	Description: "Delete a profile stored in the database",
	Res:         models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctxPkg.Get(c)
		err := ctx.DBConn.DeleteByID(&models.Profile{}, ctx.Profile.ID)
		if err != nil {
			return err
		}

		// Invalidate profiles cache
		ctx.ResetMatcherProfilesCache()

		return c.JSON(ctx.Profile)
	},
}
