package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/models"
)

func controllerRoutes(base fiber.Router) {
	control := base.Group(`/control`, requiresAuth(models.ApiKeyRoleAdmin|models.ApiKeyRoleController))
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
