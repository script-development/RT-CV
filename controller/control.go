package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/models"
)

func routeControlReloadProfiles(c *fiber.Ctx) error {
	newProfiles, err := models.GetProfiles(ctx.GetDbConn(c))
	if err != nil {
		return err
	}
	*ctx.GetProfiles(c) = newProfiles

	return c.JSON(IMap{"status": "ok"})
}
