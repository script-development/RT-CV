package controller

import (
	"errors"
	"time"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		// Debug flag can only be set by the dashboard rule
		if body.Debug && !key.Roles.ContainsSome(models.APIKeyRoleDashboard) {
			return ErrorRes(
				c,
				fiber.StatusForbidden,
				errors.New("you are not allowed to set the debug field, only api keys with the Dashboard role can set it"),
			)
		}

		err = body.CV.Validate()
		if err != nil {
			return ErrorRes(
				c,
				fiber.StatusBadRequest,
				err,
			)
		}

		// Check if we have scanned this cv earlier
		// If so return an error
		alreadyParsed, err := models.ReferenceNrAlreadyParsed(dbConn, key.ID, body.CV.ReferenceNumber)
		if alreadyParsed {
			return ErrorRes(
				c,
				fiber.StatusBadRequest,
				errors.New("a CV with this referenceNumber was previousely uploaded and parsed"),
			)
		}
		if err != nil {
			logger.WithError(err).Error("unable detect if reference number was already matched")
		}

		// Get the profiles we can use for matching
		// If they are not cached yet or the cache it outdated, set the cache
		matcherProfilesCache := ctx.GetMatcherProfilesCache(c)
		profiles := matcherProfilesCache.Profiles
		if profiles == nil || matcherProfilesCache.InsertionTime.Add(time.Hour).Before(time.Now()) {
			logger.Info("updating the profiles cache")
			// Update the cache
			profiles, err = models.GetActualActiveProfiles(dbConn)
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

		// The below is inside a goroutine to prevent blocking the fiber request
		//
		// Note that this might cause issues with slow servers when you spam the server with CV requests the go routines
		// below will not complete in time before the next request stats and thus stacking goroutines filling up the server resources
		// that could lead to 100% cpu usage or a out of memory panic
		//
		// TODO spawning a go routine seems slow
		go processMatches(processMatchesArgs{
			debug:           body.Debug,
			matchedProfiles: matchedProfiles,
			cv:              body.CV,
			logger:          *logger,
			dbConn:          dbConn,
			keyID:           key.ID,
			requestID:       requestID,
		})

		if body.Debug {
			return c.JSON(RouteScraperScanCVRes{Success: true, Matches: matchedProfiles})
		}
		return c.JSON(RouteScraperScanCVRes{Success: true})
	},
}

type processMatchesArgs struct {
	debug            bool
	matchedProfiles  []match.FoundMatch
	cv               models.CV
	logger           log.Entry
	dbConn           db.Connection
	keyID, requestID primitive.ObjectID
}

// processMatches processes the matches made to a CV
// - notify the dashboard /events page about the new match
// - label the cv reference as scanned
// - upload analytics data
// - send emails with the matches or send http requests
func processMatches(args processMatchesArgs) {
	err := dashboardListeners.publish("recived_cv", &args.requestID, args.cv)
	if err != nil {
		args.logger.WithError(err).Error("unable to save CV reference to database")
	}

	err = models.InsertParsedCVReference(args.dbConn, args.keyID, args.cv.ReferenceNumber)
	if err != nil {
		args.logger.WithError(err).Error("unable to save CV reference to database")
	}

	foundMatches := len(args.matchedProfiles) != 0

	// Insert analytics data
	if foundMatches {
		err = dashboardListeners.publish("recived_cv_matches", &args.requestID, args.matchedProfiles)
		if err != nil {
			args.logger.WithError(err).Error("unable to publish recived_cv_matches event")
		} else {
			analyticsData := make([]db.Entry, len(args.matchedProfiles))
			for idx := range args.matchedProfiles {
				args.matchedProfiles[idx].Matches.RequestID = args.requestID
				args.matchedProfiles[idx].Matches.KeyID = args.keyID
				args.matchedProfiles[idx].Matches.Debug = args.debug
				args.matchedProfiles[idx].Matches.ReferenceNr = args.cv.ReferenceNumber

				analyticsData[idx] = &args.matchedProfiles[idx].Matches
			}

			err := args.dbConn.Insert(analyticsData...)
			if err != nil {
				args.logger.WithField("analytics_entries_count", len(analyticsData)).WithError(err).Error("analytics data insertion failed")
			}
		}
	}

	if !args.debug && foundMatches {
		var pdfBytes []byte
		for _, aMatch := range args.matchedProfiles {
			if len(aMatch.Profile.OnMatch.SendMail) > 0 && pdfBytes == nil {
				// Only once and if we really need it create the email attachment pdf as this takes quite a bit of time
				//
				// MAYBE TODO:
				// Generate a pdf with placeholder values and replace the value inside the output pdf.
				// If that's possible we can speedup the pdf creation by a shitload
				pdfBytes, err = args.cv.GetPDF()
				if err != nil {
					log.WithError(err).Error("mail attachment creation error")
					return
				}
			}

			aMatch.HandleMatch(args.cv, pdfBytes)
		}
	}
}
