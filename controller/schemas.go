package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/schema"
)

// FIXME add tests for this route

func routeGetCvSchema(c *fiber.Ctx) error {
	resSchema, err := schema.From(RouteScraperScanCVBody{}, "/api/v1/schema/cv")
	if err != nil {
		return err
	}
	return c.JSON(resSchema)
}
