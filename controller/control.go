package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/models"
)

func controllerRoutes(base fiber.Router) {
	control := base.Group(`/control`, requiresAuth(models.ApiKeyRoleAdmin|models.ApiKeyRoleController))
	control.Get("/reloadProfiles", func(c *fiber.Ctx) error {
		newProfiles, err := models.GetProfiles(GetDbConn(c))
		if err != nil {
			return err
		}
		*GetProfiles(c) = newProfiles

		return c.SendString("OK")
	})
}
