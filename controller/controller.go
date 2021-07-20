package controller

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/models"
)

func Routes(app *fiber.App) {
	v1 := app.Group(`/v1`, InsertData())

	// Scraper routes
	scraper := v1.Group(`/scraper`)
	scraper.Post("/scanCV", func(c *fiber.Ctx) error {
		body := models.Cv{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles := c.UserContext().Value(Profiles(0)).(*[]models.Profile)
		match.Match("werk.nl", *profiles, body)

		return c.SendString("Ok")
	})

	// Control routes
	control := v1.Group(`/control`)
	control.Get("/reloadMatches", func(c *fiber.Ctx) error {
		profiles := c.UserContext().Value(Profiles(0)).(*[]models.Profile)

		newProfiles, err := models.GetProfiles()
		if err != nil {
			return err
		}
		*profiles = newProfiles

		return c.SendString("Ok")
	})
}

type Profiles uint8

// InsertData adds the profiles to every route
func InsertData() fiber.Handler {
	profiles, err := models.GetProfiles()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.WithValue(context.Background(), Profiles(0), &profiles)

	return func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	}
}
