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

		MatchesProcess.AppendMatchesToProcess(ProcessMatches{
			Debug:           body.Debug,
			MatchedProfiles: matchedProfiles,
			CV:              body.CV,
			Logger:          *logger,
			DBConn:          dbConn,
			KeyID:           key.ID,
			RequestID:       requestID,
		})

		if body.Debug {
			return c.JSON(RouteScraperScanCVRes{Success: true, Matches: matchedProfiles})
		}
		return c.JSON(RouteScraperScanCVRes{Success: true})
	},
}

// MatchesProcessor is a struct that contains a list of matches to be processed in the background
//
// To register a match to be processed call (*MatchesProcessor).AppendMatchesToProcess
// After calling the above (*MatchesProcessor).processMatches should automatically pick up the match and handle it
type MatchesProcessor struct {
	list    []ProcessMatches
	c       *sync.Cond
	started bool
}

// MatchesProcess holts the process matches that should be processed in the background
var MatchesProcess = &MatchesProcessor{
	list:    []ProcessMatches{},
	c:       sync.NewCond(&sync.Mutex{}),
	started: false,
}

// AppendMatchesToProcess adds a list of matches to be processed by (*MatchesProcessor).processMatches
func (p *MatchesProcessor) AppendMatchesToProcess(args ProcessMatches) {
	p.c.L.Lock()
	p.list = append(p.list, args)
	p.c.Signal()
	if !p.started {
		p.started = true
		go p.processMatches()
	}
	p.c.L.Unlock()
}

// processMatches is a process that should be running in the background that process matches
func (p *MatchesProcessor) processMatches() {
	for {
		p.c.L.Lock()
		if len(p.list) == 0 {
			// Once the list is empty wait for a signal for the list to fill up again
			p.c.Wait()
		}
		matchToProcess := p.list[0]
		p.list = p.list[1:]
		p.c.L.Unlock()
		matchToProcess.Process()
	}
}

// ProcessMatches contains the content for processing a match
type ProcessMatches struct {
	Debug            bool
	MatchedProfiles  []match.FoundMatch
	CV               models.CV
	Logger           log.Entry
	DBConn           db.Connection
	KeyID, RequestID primitive.ObjectID
}

// Process processes the matches made to a CV
// - notify the dashboard /events page about the new match
// - safe the matches of this reference number for analytics and for detecting duplicates
// - send emails with the matches or send http requests
func (args ProcessMatches) Process() {
	err := dashboardListeners.publish("recived_cv", &args.RequestID, args.CV)
	if err != nil {
		args.Logger.WithError(err).Error("unable to save CV reference to database")
	}

	if len(args.MatchedProfiles) == 0 {
		return
	}

	// Get earlier matches on this reference number
	earlierMatches, err := models.GetMatchesOnReferenceNr(args.DBConn, args.CV.ReferenceNumber, &args.KeyID)
	if err != nil {
		args.Logger.WithError(err).Error("unable to execute query to get earlier made matches to this reference number")
		earlierMatches = []models.Match{}
	}

	// Remove matches that where already made earlier
	// We loop in reverse so we can remove items from the slice
	for idx := len(args.MatchedProfiles) - 1; idx >= 0; idx-- {
		for _, earlierMatche := range earlierMatches {
			if args.MatchedProfiles[idx].Profile.ID == earlierMatche.ProfileID {
				args.MatchedProfiles = append(args.MatchedProfiles[:idx], args.MatchedProfiles[idx+1:]...)
				break
			}
		}
	}

	// Re-check the amount of matched profiles as we might have filtered out at the step above
	if len(args.MatchedProfiles) == 0 {
		return
	}

	err = dashboardListeners.publish("recived_cv_matches", &args.RequestID, args.MatchedProfiles)
	if err != nil {
		args.Logger.WithError(err).Error("unable to publish recived_cv_matches event")
	}

	analyticsData := make([]db.Entry, len(args.MatchedProfiles))
	for idx := range args.MatchedProfiles {
		args.MatchedProfiles[idx].Matches.RequestID = args.RequestID
		args.MatchedProfiles[idx].Matches.KeyID = args.KeyID
		args.MatchedProfiles[idx].Matches.Debug = args.Debug
		args.MatchedProfiles[idx].Matches.ReferenceNr = args.CV.ReferenceNumber

		analyticsData[idx] = &args.MatchedProfiles[idx].Matches
	}
	err = args.DBConn.Insert(analyticsData...)
	if err != nil {
		args.Logger.WithField("analytics_entries_count", len(analyticsData)).WithError(err).Error("analytics data insertion failed")
	}

	if args.Debug {
		return
	}

	var pdfBytes []byte
	for _, aMatch := range args.MatchedProfiles {
		if len(aMatch.Profile.OnMatch.SendMail) > 0 && pdfBytes == nil {
			// Only once and if we really need it create the email attachment pdf as this takes quite a bit of time
			//
			// MAYBE TODO:
			// Generate a pdf with placeholder values and replace the value inside the output pdf.
			// If that's possible we can speedup the pdf creation by a shitload
			pdfBytes, err = args.CV.GetPDF()
			if err != nil {
				log.WithError(err).Error("mail attachment creation error")
				return
			}
		}

		aMatch.HandleMatch(args.CV, pdfBytes)
	}
}
