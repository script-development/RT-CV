package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	reqPkg "github.com/script-development/RT-CV/controller/ctx"
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
	Success    bool               `json:"success"`
	HasMatches bool               `json:"hasMatches"`
	Matches    []match.FoundMatch `json:"matches" jsonSchema:"hidden" description:"Only contains matches if the debug property is set to true"`
}

var routeScraperScanCV = routeBuilder.R{
	Description: "Main route to scrape the CV",
	Res:         RouteScraperScanCVRes{},
	Body:        RouteScraperScanCVBody{},
	Fn: func(c *fiber.Ctx) error {
		ctx := reqPkg.Get(c)

		body := RouteScraperScanCVBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		// Debug flag can only be set by the dashboard rule
		if body.Debug && !ctx.Key.Roles.ContainsSome(models.APIKeyRoleDashboard) {
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
		profiles := ctx.MatcherProfilesCache.Profiles
		if profiles == nil || ctx.MatcherProfilesCache.InsertionTime.Add(time.Hour*24).Before(time.Now()) {
			ctx.Logger.Info("updating the profiles cache")
			// Update the cache
			profilesFromDB, err := models.GetActualActiveProfiles(ctx.DBConn)
			if err != nil {
				return err
			}
			profiles = make([]*models.Profile, len(profilesFromDB))
			for idx := range profilesFromDB {
				profiles[idx] = &profilesFromDB[idx]
			}
			*ctx.MatcherProfilesCache = reqPkg.MatcherProfilesCache{
				Profiles:      profiles,
				InsertionTime: time.Now(),
			}
		}

		// Try to match a profile to a CV
		matchedProfiles := match.Match(ctx.Key.ID, ctx.RequestID, profiles, body.CV)

		resp := RouteScraperScanCVRes{Success: true, Matches: []match.FoundMatch{}}
		if len(matchedProfiles) == 0 {
			return c.JSON(resp)
		}

		MatchesProcess.AppendMatchesToProcess(ProcessMatches{
			Debug:           body.Debug,
			MatchedProfiles: matchedProfiles,
			CV:              body.CV,
			Logger:          *ctx.Logger,
			DBConn:          ctx.DBConn,
			KeyID:           ctx.Key.ID,
			KeyName:         ctx.Key.Name,
		})

		resp.HasMatches = true
		if body.Debug {
			resp.Matches = matchedProfiles
		}

		return c.JSON(resp)
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
	Debug           bool
	MatchedProfiles []match.FoundMatch
	CV              models.CV
	Logger          log.Entry
	DBConn          db.Connection
	KeyID           primitive.ObjectID
	KeyName         string
}

// DataSendToHook contains the content for processing a match
type DataSendToHook struct {
	MatchedProfiles []match.FoundMatch `json:"matchedProfiles"`
	CV              models.CV          `json:"cv"`
	KeyID           primitive.ObjectID `json:"keyId" description:"The ID of the API key that was used to upload this CV"`
	KeyName         string             `json:"keyName" description:"The Name of the API key that was used to upload this CV"`
	IsTest          bool               `json:"isTest" description:"True if this hook call was manually triggered"`
}

// Process processes the matches made to a CV
// - notify the dashboard /events page about the new match
// - safe the matches of this reference number for analytics and for detecting duplicates
// - send emails with the matches or send http requests
func (args ProcessMatches) Process() {
	if len(args.MatchedProfiles) == 0 {
		return
	}

	// Re-check the amount of matched profiles as we might have filtered out at the step above
	if len(args.MatchedProfiles) == 0 {
		return
	}

	hooks := []models.OnMatchHook{}
	err := args.DBConn.Find(&models.OnMatchHook{}, &hooks, nil)
	if err != nil {
		args.Logger.WithError(err).Error("Finding on match hooks failed")
		return
	}

	if len(hooks) == 0 {
		log.Error("no on match hooks configured")
		return
	}

	hookData, err := json.Marshal(DataSendToHook{
		MatchedProfiles: args.MatchedProfiles,
		CV:              args.CV,
		KeyID:           args.KeyID,
		KeyName:         args.KeyName,
	})
	if err != nil {
		args.Logger.WithError(err).Error("creating hook data failed")
		return
	}

	for _, hook := range hooks {
		if err != nil {
			args.Logger.WithError(err).Error("creating hook data failed")
			continue
		}

		err = hook.Call(bytes.NewBuffer(hookData))
		if err != nil {
			args.Logger.WithError(err).Error("creating hook data failed")
		} else {
			args.Logger.WithField("hook", hook.URL).WithField("hook_id", hook.ID.Hex()).Info("hook called")
		}
	}
}
