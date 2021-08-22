package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
)

func routeAuthSeed(c *fiber.Ctx) error {
	return c.JSON(IMap{
		"seed": string(ctx.GetAuth(c).GetBaseSeed()),
	})
}
