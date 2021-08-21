package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/models"
)

func routeScraperScanCV(c *fiber.Ctx) error {
	body := models.CV{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	profiles := ctx.GetProfiles(c)
	matchedProfiles := match.Match("werk.nl", *profiles, body) // TODO remove this hardcoded value
	if len(matchedProfiles) > 0 {
		for _, profile := range matchedProfiles {
			_, err := body.GetPDF(profile, "") // TODO add matchtext
			if err != nil {
				return fmt.Errorf("unable to generate PDF from CV, err: %s", err.Error())
			}
			fmt.Println(profile.Emails)
			// for _, email := range profile.Emails {
			// 	email.Email.Name
			// }
		}
	}

	return c.SendString("OK")
}
