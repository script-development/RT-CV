package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/models"
)

func Routes(app *fiber.App) {
	v1 := app.Group(`/v1`, InsertData())

	auth := v1.Group(`/auth`)
	auth.Get(`/salt`, func(c *fiber.Ctx) error {
		// TODO: Implement me
		return nil
	})

	// Scraper routes
	scraper := v1.Group(`/scraper`, requiresAuth(models.ApiKeyRoleScraper))
	scraper.Post("/scanCV", func(c *fiber.Ctx) error {
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
				// for _, email := range profile.Emails {
				// 	email.Email.Name
				// }
			}
		}

		return c.SendString("OK")
	})

	// Control routes
	control := v1.Group(`/control`, requiresAuth(models.ApiKeyRoleAdmin|models.ApiKeyRoleController))
	control.Get("/reloadProfiles", func(c *fiber.Ctx) error {
		profiles := c.UserContext().Value(ProfilesCtxKey).(*[]models.Profile)

		newProfiles, err := models.GetProfiles()
		if err != nil {
			return err
		}
		*profiles = newProfiles

		return c.SendString("OK")
	})
}

type ProfilesCtx uint8
type AuthCtx uint8
type KeyCtx uint8
type LoggerCtx uint8

const (
	ProfilesCtxKey = ProfilesCtx(0)
	AuthCtxKey     = AuthCtx(0)
	KeyCtxKey      = KeyCtx(0)
	LoggerCtxKey   = LoggerCtx(0)
)
