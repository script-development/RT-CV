package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/models"
)

func controllerRoutes(base fiber.Router) {
	control := base.Group(`/control`, requiresAuth(models.ApiKeyRoleAdmin|models.ApiKeyRoleController))
	control.Get("/reloadProfiles", routeControlReloadProfiles)
}

func routeControlReloadProfiles(c *fiber.Ctx) error {
	newProfiles, err := models.GetProfiles(GetDbConn(c))
	if err != nil {
		return err
	}
	*GetProfiles(c) = newProfiles

	return c.JSON(map[string]string{"status": "ok"})
}
