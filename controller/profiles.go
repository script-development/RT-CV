package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeAllProfiles = routeBuilder.R{
	Description: "get all profiles stored in the database",
	Res:         []models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		profiles, err := models.GetProfiles(ctx.GetDbConn(c))
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
		db := ctx.GetDbConn(c)

		profilesCount, err := models.GetProfilesCount(db)
		if err != nil {
			return err
		}

		usableProfilesCount, err := models.GetActualActiveProfilesCount(db)
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

			dbConn := ctx.GetDbConn(c)
			profile, err := models.GetProfile(dbConn, profileID)
			if err != nil {
				return err
			}

			c.SetUserContext(
				ctx.SetProfile(
					c.UserContext(),
					&profile,
				),
			)
			return c.Next()
		},
	}
}

var routeGetProfile = routeBuilder.R{
	Description: "get a profile based on it's ID from the database",
	Res:         models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		profile := ctx.GetProfile(c)
		return c.JSON(profile)
	},
}

var routeCreateProfile = routeBuilder.R{
	Description: "create a new profile that can match scraped CVs",
	Res:         models.Profile{},
	Body:        models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		var profile models.Profile
		err := c.BodyParser(&profile)
		if err != nil {
			return err
		}

		err = profile.ValidateCreateNewProfile()
		if err != nil {
			return err
		}

		// Set the ID of the profile
		profile.M = db.NewM()

		// Save the profile to the database
		dbConn := ctx.GetDbConn(c)
		err = dbConn.Insert(&profile)
		if err != nil {
			return err
		}

		// Invalidate profiles cache
		*ctx.GetMatcherProfilesCache(c) = ctx.MatcherProfilesCache{}

		return c.JSON(profile)
	},
}

func routeModifyProfile(c *fiber.Ctx) error {
	profile := ctx.GetProfile(c)
	// FIXME implement route

	// Invalidate profiles cache
	*ctx.GetMatcherProfilesCache(c) = ctx.MatcherProfilesCache{}

	return c.JSON(profile)
}

var routeDeleteProfile = routeBuilder.R{
	Description: "Delete a profile stored in the database",
	Res:         models.Profile{},
	Fn: func(c *fiber.Ctx) error {
		profile := ctx.GetProfile(c)
		dbConn := ctx.GetDbConn(c)
		err := dbConn.DeleteByID(profile)
		if err != nil {
			return err
		}

		// Invalidate profiles cache
		*ctx.GetMatcherProfilesCache(c) = ctx.MatcherProfilesCache{}

		return c.JSON(profile)
	},
}
