package controller

// This file contains json schemas for certain types
// These are used by the dashboard to to validate user input data
//
// To be specific the monaco editor (the code editor we use in the dashboard) uses json schemas to validate the user input

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/schema"
)

func routeGetCvSchema(c *fiber.Ctx) error {
	resSchema, err := schema.From(RouteScraperScanCVBody{}, "/api/v1/schema/cv")
	if err != nil {
		return err
	}
	return c.JSON(resSchema)
}
