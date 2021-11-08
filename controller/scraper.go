package controller

import (
	"errors"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

// RouteScraperScanCVBody is the request body of the routeScraperScanCV
type RouteScraperScanCVBody struct {
	CV    models.CV `json:"cv"`
	Debug bool      `json:"debug" jsonSchema:"hidden"`
}

// RouteScraperScanCVRes contains the response data of routeScraperScanCV
type RouteScraperScanCVRes struct {
	Success bool `json:"success"`

	// Matches is only set if the debug property is set
	Matches []match.FoundMatch `json:"matches" jsonSchema:"hidden"`
}

var routeScraperScanCV = routeBuilder.R{
	Description: "Main route to scrape the CV",
	Res:         RouteScraperScanCVRes{},
	Body:        RouteScraperScanCVBody{},
	Fn: func(c *fiber.Ctx) error {
		key := ctx.GetKey(c)
		requestID := ctx.GetRequestID(c)
		dbConn := ctx.GetDbConn(c)
		logger := ctx.GetLogger(c)

		body := RouteScraperScanCVBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		err = body.CV.Validate()
		if err != nil {
			return ErrorRes(
				c,
				fiber.StatusBadRequest,
				err,
			)
		}

		alreadyParsed, err := models.ReferenceNrAlreadyParsed(dbConn, key.ID, body.CV.ReferenceNumber)
		if alreadyParsed {
			return ErrorRes(c, fiber.StatusBadRequest, errors.New("a CV with this referenceNumber was previousely uploaded and parsed"))
		}
		if err != nil {
			logger.WithError(err).Error("unable detect if reference number was already matched")
		}

		err = models.InsertParsedCVReference(dbConn, key.ID, body.CV.ReferenceNumber)
		if err != nil {
			logger.WithError(err).Error("unable to save CV reference to database")
		}

		err = dashboardListeners.publish("recived_cv", &requestID, body.CV)
		if err != nil {
			return err
		}

		if body.Debug && !key.Roles.ContainsSome(models.APIKeyRoleDashboard) {
			return ErrorRes(
				c,
				fiber.StatusForbidden,
				errors.New("you are not allowed to set the debug field, only api keys with the Dashboard role can set it"),
			)
		}

		matcherProfilesCache := ctx.GetMatcherProfilesCache(c)
		profiles := matcherProfilesCache.Profiles
		if profiles == nil || matcherProfilesCache.InsertionTime.Add(time.Hour).Before(time.Now()) {
			// Update the cache
			profiles, err = models.GetActiveProfilesWithOnMatch(dbConn)
			if err != nil {
				return err
			}
			*matcherProfilesCache = ctx.MatcherProfilesCache{
				Profiles:      profiles,
				InsertionTime: time.Now(),
			}
		}

		// Try to match a profile to a CV
		matchedProfiles := match.Match(key.Domains, profiles, body.CV)
		foundMatches := len(matchedProfiles) != 0

		// Insert analytics data
		if foundMatches {
			err = dashboardListeners.publish("recived_cv_matches", &requestID, matchedProfiles)
			if err != nil {
				return err
			}

			analyticsData := make([]db.Entry, len(matchedProfiles))
			for idx := range matchedProfiles {
				matchedProfiles[idx].Matches.RequestID = requestID
				matchedProfiles[idx].Matches.KeyID = key.ID
				matchedProfiles[idx].Matches.Debug = body.Debug

				analyticsData[idx] = &matchedProfiles[idx].Matches
			}

			go func(logger *log.Entry, analyticsData []db.Entry) {
				err := dbConn.Insert(analyticsData...)
				if err != nil {
					logger.WithError(err).Error("analytics data insertion failed")
				}
			}(logger.WithField("analytics_entries_count", len(analyticsData)), analyticsData)
		}

		if body.Debug {
			return c.JSON(RouteScraperScanCVRes{Success: true, Matches: matchedProfiles})
		}

		if foundMatches {
			logger.Infof("found %d matches", len(matchedProfiles))

			var wg sync.WaitGroup

			for _, aMatch := range matchedProfiles {
				err := aMatch.HandleMatch(&wg, &body.CV)
				if err != nil {
					log.WithError(err).Error("sending match error")
				}
			}

			wg.Wait()
		}

		return c.JSON(RouteScraperScanCVRes{Success: true})
	},
}
