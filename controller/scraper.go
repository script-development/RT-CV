package controller

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

// RouteScraperScanCVBody is the request body of the routeScraperScanCV
type RouteScraperScanCVBody struct {
	CV    models.CV `json:"cv"`
	Debug bool      `json:"debug"`
}

// TODO: maybe we should not return the actual profiles matched, this exposes information not meant for this api key user
var routeScraperScanCV = routeBuilder.R{
	Description: "Main route to scrape the CV",
	Res:         []match.AMatch{},
	Body:        RouteScraperScanCVBody{},
	Fn: func(c *fiber.Ctx) error {
		key := ctx.GetKey(c)

		body := RouteScraperScanCVBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles := ctx.GetProfiles(c)
		matchedProfiles := match.Match(key.Domains, *profiles, body.CV)
		if body.Debug {
			return c.JSON(matchedProfiles)
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

		return c.JSON(matchedProfiles)
	},
}
