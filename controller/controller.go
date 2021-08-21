package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/models"
)

// Group is a wrapper around fiber.Router.Group to change it into a callback function with the group
// This makes it easier to understand the routes structure
func Group(c fiber.Router, prefix string, group func(fiber.Router), handlers ...func(*fiber.Ctx) error) {
	group(c.Group(prefix, handlers...))
}

func Routes(app *fiber.App, dbConn dbInterfaces.Connection, serverSeed []byte) {
	Group(app, `/v1`, func(c fiber.Router) {
		Group(c, `/auth`, func(c fiber.Router) {
			c.Get(`/seed`, routeAuthSeed)
		})

		Group(c, `/scraper`, func(c fiber.Router) {
			c.Post(`/scanCV`, routeScraperScanCV)
			Group(c, `/secret/:key`, func(c fiber.Router) {
				c.Delete(``, routeDeleteSecret)
				Group(c, `/:secretKey`, func(c fiber.Router) {
					c.Post(``, routeCreateSecret)
					c.Put(``, routeUpdateSecret)
					c.Get(``, routeGetSecret)
				}, validSecretKeyMiddleware())
			}, validKeyMiddleware())
		}, requiresAuth(models.ApiKeyRoleScraper))

		Group(c, `/control`, func(c fiber.Router) {
			c.Get(`/reloadProfiles`, routeControlReloadProfiles)
		}, requiresAuth(models.ApiKeyRoleAdmin|models.ApiKeyRoleController))
	}, InsertData(dbConn, serverSeed))
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	return ErrorRes(c, 500, err)
}

func ErrorRes(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(map[string]string{
		"error": err.Error(),
	})
}
