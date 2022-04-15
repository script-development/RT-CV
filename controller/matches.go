package controller

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	ctxPkg "github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const fromToExplainer = `
the from and to param should be in the RFC 3339 format
_RFC 3339 is basically an extension to iso 8601_`

var routeGetMatchesPeriod = routeBuilder.R{
	Description: `get all matches made within a certain period.` + fromToExplainer,
	Res:         []models.Match{},
	Fn:          func(c *fiber.Ctx) error { return getMatchesForPeriod(c, nil) },
}

var routeGetMatchesPeroidForProfile = routeBuilder.R{
	Description: `get all matches made within a certain period for a spesific profile.` + fromToExplainer,
	Res:         []models.Match{},
	Fn:          func(c *fiber.Ctx) error { return getMatchesForPeriod(c, bson.M{"profileId": ctxPkg.Get(c).Profile.ID}) },
}

var routeGetMatchesPeroidForProfilesQuery = routeBuilder.R{
	Description: `get all matches made within a certain period for a spesific range of profiles based on the MongoDB query defined in the request body.` + fromToExplainer + `

For more info about mongodb queries you can take a look at: https://www.mongodb.com/docs/manual/tutorial/query-documents/#std-label-read-operations-query-argument
Note that you can't filter for the _id field`,
	Body: bson.M{},
	Res:  []models.Match{},
	Fn: func(c *fiber.Ctx) error {
		body := bson.M{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles, err := models.GetProfiles(ctxPkg.Get(c).DBConn, body)
		if err != nil {
			return err
		}

		profileIDs := make([]primitive.ObjectID, len(profiles))
		for idx, profile := range profiles {
			profileIDs[idx] = profile.ID
		}

		return getMatchesForPeriod(c, bson.M{"profileId": bson.M{"$in": profileIDs}})
	},
}

func getMatchesForPeriod(c *fiber.Ctx, extraQueryFilters bson.M) error {
	ctx := ctxPkg.Get(c)

	fromParam, toParam := c.Params("from"), c.Params("to")
	from, err := time.Parse(time.RFC3339, fromParam)
	if err != nil {
		return errors.New("unable \"from\" parse from param as RFC 3339")
	}
	to, err := time.Parse(time.RFC3339, toParam)
	if err != nil {
		return errors.New("unable \"to\" parse from param as RFC 3339")
	}

	now := time.Now()
	if to.Before(now) {
		// We cannot do things in the past, we're fine with caching

		// 1 second * 60 = minute * 60 = hour * 24 = day * 14 = 2 weeks
		maxAge := 1 * 60 * 60 * 24 * 14
		c.Response().Header.Set("Cache-Control", "private")
		c.Response().Header.Add("Cache-Control", "max-age="+strconv.Itoa(maxAge))

		c.Response().Header.SetLastModified(to)
	}

	query := extraQueryFilters
	if query == nil {
		query = bson.M{}
	}

	query["when"] = bson.M{
		"$gt": from,
		"$lt": to,
	}

	matches := []models.Match{}
	err = ctx.DBConn.Find(&models.Match{}, &matches, query)
	if err != nil {
		return err
	}

	return c.JSON(matches)
}

// MatchesPerProfile contains a map of profile IDs and their amound of matches
type MatchesPerProfile map[primitive.ObjectID]ProfileMatches

// ProfileMatches conatins the amound of matches for a spesific profile
type ProfileMatches struct {
	Unique uint64 `json:"unique"`
	Total  uint64 `json:"total"`
}

var routeGetMatchesPeriodPerProfile = routeBuilder.R{
	Description: `get all matches made within a certain period based on the profiles.
The map key contains the profile id and the value contains the amound of matches.
The from and to param should be in the RFC 3339 format
_RFC 3339 is basically an extension to iso 8601_`,
	Res: MatchesPerProfile{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctxPkg.Get(c)

		fromParam, toParam := c.Params("from"), c.Params("to")
		from, err := time.Parse(time.RFC3339, fromParam)
		if err != nil {
			return errors.New("unable \"from\" parse from param as RFC 3339")
		}
		to, err := time.Parse(time.RFC3339, toParam)
		if err != nil {
			return errors.New("unable \"to\" parse from param as RFC 3339")
		}

		now := time.Now()
		if to.Before(now) {
			// We cannot do things in the past, we're fine with caching

			// 1 second * 60 = minute * 60 = hour * 24 = day * 14 = 2 weeks
			maxAge := 1 * 60 * 60 * 24 * 14
			c.Response().Header.Set("Cache-Control", "private")
			c.Response().Header.Add("Cache-Control", "max-age="+strconv.Itoa(maxAge))

			c.Response().Header.SetLastModified(to)
		}

		query := bson.M{
			"when": bson.M{
				"$gt": from,
				"$lt": to,
			},
		}

		if c.Params(`profile`) != "" {
			query["profileId"] = ctx.Profile.ID
		}

		matches := []models.Match{}
		err = ctx.DBConn.Find(&models.Match{}, &matches, query)
		if err != nil {
			return err
		}

		// First map key = id of profile
		// Second map key = reference number
		// Seoncd map value = amound of times matched this reference number
		matchesMap := map[primitive.ObjectID]map[string]uint64{}
		for _, match := range matches {
			matchesOnProfile, ok := matchesMap[match.ProfileID]
			if !ok {
				matchesOnProfile = map[string]uint64{}
			}
			matchesOnProfile[match.ReferenceNr]++
			matchesMap[match.ProfileID] = matchesOnProfile
		}

		resp := MatchesPerProfile{}
		for profileID, matchesOnProfile := range matchesMap {
			matchesCount := ProfileMatches{}
			for _, numberOfMatches := range matchesOnProfile {
				matchesCount.Unique++
				matchesCount.Total += numberOfMatches
			}
			resp[profileID] = matchesCount
		}

		return c.JSON(resp)
	},
}
