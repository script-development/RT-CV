package controller

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/schema"
)

// FIXME add tests for this route

func routeGetCvSchema(c *fiber.Ctx) error {
	properties, requiredFields, err := schema.From(reflect.TypeOf(RouteScraperScanCVBody{}))
	if err != nil {
		return err
	}
	// FIXME move the code below to schema.From
	res := schema.Property{
		Schema:     schema.VersionUsed,
		ID:         "/api/v1/schema/cv",
		Type:       schema.PropertyTypeObject,
		Properties: properties,
		Required:   requiredFields,
	}
	return c.JSON(res)
}
