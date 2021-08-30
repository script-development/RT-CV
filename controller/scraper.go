package controller

import (
	"fmt"

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
			for _, aMatch := range matchedProfiles {
				_, err := body.CV.GetPDF(aMatch.Profile, aMatch.Matches.GetMatchSentence())
				if err != nil {
					return fmt.Errorf("unable to generate PDF from CV, err: %s", err.Error())
				}
				fmt.Println(aMatch.Profile.Emails)
				// for _, email := range profile.Emails {
				// 	email.Email.Name
				// }
			}
		}

		return c.JSON(matchedProfiles)
	},
}
