package controller

import (
	"errors"
	"fmt"
	"sync"

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

// TODO: maybe we should not return the actual profiles matched, this exposes information not meant for this api key user
var routeScraperScanCV = routeBuilder.R{
	Description: "Main route to scrape the CV",
	Res:         RouteScraperScanCVRes{},
	Body:        RouteScraperScanCVBody{},
	Fn: func(c *fiber.Ctx) error {
		key := ctx.GetKey(c)
		requestID := ctx.GetRequestID(c)
		dbConn := ctx.GetDbConn(c)

		body := RouteScraperScanCVBody{}
		err := c.BodyParser(&body)
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

		profiles := ctx.GetProfiles(c)
		matchedProfiles := match.Match(key.Domains, *profiles, body.CV)

		// Insert analytics data
		analyticsData := make([]db.Entry, len(matchedProfiles))
		for idx := range matchedProfiles {
			matchedProfiles[idx].Matches.RequestID = requestID
			matchedProfiles[idx].Matches.KeyID = key.ID
			matchedProfiles[idx].Matches.Debug = body.Debug

			analyticsData[idx] = &matchedProfiles[idx].Matches
		}
		go dbConn.Insert(analyticsData...)

		if body.Debug {
			return c.JSON(RouteScraperScanCVRes{Success: true, Matches: matchedProfiles})
		}

		if len(matchedProfiles) > 0 {
			var wg sync.WaitGroup

			for _, aMatch := range matchedProfiles {
				_, err := body.CV.GetPDF(aMatch.Profile, aMatch.Matches.GetMatchSentence())
				if err != nil {
					return fmt.Errorf("unable to generate PDF from CV, err: %s", err.Error())
				}

				wg.Add(len(aMatch.Profile.OnMatch.HTTPCall) + len(aMatch.Profile.OnMatch.SendMail))

				for _, http := range aMatch.Profile.OnMatch.HTTPCall {
					go func(http models.ProfileHTTPCallData) {
						http.MakeRequest()
						wg.Done()
					}(http)
				}

				for _, email := range aMatch.Profile.OnMatch.SendMail {
					go func(email models.ProfileSendEmailData) {
						email.SendEmail()
						wg.Done()
					}(email)
				}
			}

			wg.Wait()
		}

		return c.JSON(RouteScraperScanCVRes{Success: true})
	},
}
