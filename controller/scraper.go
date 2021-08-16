package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/models"
)

func scraperRoutes(base fiber.Router) {
	scraper := base.Group(`/scraper`, requiresAuth(models.ApiKeyRoleScraper))
	scraper.Post(`/scanCV`, func(c *fiber.Ctx) error {
		body := models.Cv{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles := c.UserContext().Value(ProfilesCtxKey).(*[]models.Profile)
		matchedProfiles := match.Match("werk.nl", *profiles, body) // TODO remove this hardcoded value
		if len(matchedProfiles) > 0 {
			for _, profile := range matchedProfiles {
				_, err := body.GetPDF(profile, "") // TODO add matchtext
				if err != nil {
					return fmt.Errorf("unable to generate PDF from CV, err: %s", err.Error())
				}
				fmt.Println("")
				// for _, email := range profile.Emails {
				// 	email.Email.Name
				// }
			}
		}

		return c.SendString("OK")
	})

	secret := scraper.Group(`/secret`)
	secret.Get(`/:key`, func(c *fiber.Ctx) error {
		c.Params("key")
		// TODO: Get secret
		return nil
	})
	secret.Put(`/:key`, func(c *fiber.Ctx) error {
		c.Params("key")
		// TODO: Create secret
		return nil
	})
	secret.Delete(`/:key`, func(c *fiber.Ctx) error {
		c.Params("key")
		// TODO: Delete secret
		return nil
	})
}
