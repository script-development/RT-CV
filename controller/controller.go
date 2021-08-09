package controller

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/models"
)

func Routes(app *fiber.App) {
	v1 := app.Group(`/v1`, InsertData())

	// Scraper routes
	scraper := v1.Group(`/scraper`, requiresAuth("scraper"))
	scraper.Post("/scanCV", func(c *fiber.Ctx) error {
		body := models.Cv{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles := c.UserContext().Value(Profiles(0)).(*[]models.Profile)
		matchedProfiles := match.Match("werk.nl", *profiles, body)
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
	control := v1.Group(`/control`, requiresAuth("admin"))
	control.Get("/reloadProfiles", func(c *fiber.Ctx) error {
		profiles := c.UserContext().Value(Profiles(0)).(*[]models.Profile)

		newProfiles, err := models.GetProfiles()
		if err != nil {
			return err
		}
		*profiles = newProfiles

		return c.SendString("OK")
	})
}

type Profiles uint8
type Auth uint8

// InsertData adds the profiles to every route
func InsertData() fiber.Handler {
	profiles, err := models.GetProfiles()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.WithValue(context.Background(), Profiles(0), &profiles)
	ctx = context.WithValue(ctx, Auth(0), auth.New())

	return func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func requiresAuth(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
