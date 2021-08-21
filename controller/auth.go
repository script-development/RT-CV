package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
)

func routeAuthSeed(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"seed": string(ctx.GetAuth(c).GetBaseSeed()),
	})
}
