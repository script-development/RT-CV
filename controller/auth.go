package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

var routeGetKeyInfo = routeBuilder.R{
	Description: "Get information about the key you are using to authenticate with",
	Res:         models.APIKeyInfo{},
	Fn: func(c *fiber.Ctx) error {
		return c.JSON(ctx.Get(c).Key.Info())
	},
}
