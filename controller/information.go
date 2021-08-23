package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/models"
)

func routeGetKeys(c *fiber.Ctx) error {
	dbConn := ctx.GetDbConn(c)
	keys, err := models.GetAPIKeys(dbConn)
	if err != nil {
		return err
	}
	return c.JSON(keys)
}
