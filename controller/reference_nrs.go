package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
)

var scannedReferenceNrs = routeBuilder.R{
	Description: "DEPRECATED and will always return an empty array",
	Res:         []string{},
	Fn: func(c *fiber.Ctx) error {
		return c.JSON([]string{})
	},
}
