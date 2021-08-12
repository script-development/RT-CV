package controller

import (
	"context"
	"errors"
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
	scraper := v1.Group(`/scraper`, requiresAuth(models.ApiKeyRoleScraper))
	scraper.Post("/scanCV", func(c *fiber.Ctx) error {
		body := models.Cv{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles := c.UserContext().Value(ProfilesCtxKey).(*[]models.Profile)
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

type Profiles uint8
type Auth uint8
type Salt uint8
type Roles uint8

const (
	ProfilesCtxKey = Profiles(0)
	AuthCtxKey     = Auth(0)
	SaltCtxKey     = Salt(0)
	RolesCtxKey    = Roles(0)
)

// InsertData adds the profiles to every route
func InsertData() fiber.Handler {
	profiles, err := models.GetProfiles()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.WithValue(context.Background(), ProfilesCtxKey, &profiles)

	keys, err := models.GetApiKeys()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx = context.WithValue(ctx, AuthCtxKey, auth.New(keys))

	return func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func requiresAuth(requiredRoles models.ApiKeyRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		auth := ctx.Value(AuthCtxKey).(*auth.Auth)

		authorizationHeader := []byte(c.Get("Authorization"))
		key, salt, err := auth.Authenticate(authorizationHeader)
		if err != nil {
			return err
		}

		if !key.Roles.ContainsSome(requiredRoles) {
			return errors.New("you do not have the permissions to access this route")
		}

		ctx = context.WithValue(ctx, RolesCtxKey, key.Roles)
		ctx = context.WithValue(ctx, SaltCtxKey, salt)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
