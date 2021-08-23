package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// IMap is a wrapper around map[string]interface{} that's faster to use
type IMap map[string]interface{}

// Group is a wrapper around fiber.Router.Group to change it into a callback function with the group
// This makes it easier to understand the routes structure
func Group(c fiber.Router, prefix string, group func(fiber.Router), handlers ...func(*fiber.Ctx) error) {
	group(c.Group(prefix, handlers...))
}

// Routes defines the routes used
func Routes(app *fiber.App, dbConn db.Connection, serverSeed []byte) {
	Group(app, `/api/v1`, func(c fiber.Router) {
		Group(c, `/auth`, func(c fiber.Router) {
			c.Get(`/keyinfo`, requiresAuth(0), routeGetKeyInfo)
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
		}, requiresAuth(models.APIKeyRoleScraper))

		Group(c, `/control`, func(c fiber.Router) {
			Group(c, `/profiles`, func(c fiber.Router) {
				c.Post(``, routeCreateProfile)
				c.Get(``, routeAllProfiles)
				Group(c, `/:profile`, func(c fiber.Router) {
					c.Get(``, routeGetProfile)
					// c.Put(``, routeModifyProfile) // TODO
					c.Delete(``, routeDeleteProfile)
				}, middlewareBindProfile())
			})
		}, requiresAuth(models.APIKeyRoleController))

		Group(c, `/keys`, func(c fiber.Router) {
			c.Get(``, routeGetKeys)
		}, requiresAuth(models.APIKeyRoleInformationObtainer|models.APIKeyRoleController))
	}, InsertData(dbConn, serverSeed))
}

// FiberErrorHandler handles errors in fiber
// In our case that means we change the errors from text to json
func FiberErrorHandler(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrorRes(c, 404, errors.New("item not found"))
	}
	return ErrorRes(c, 500, err)
}

// ErrorRes returns the error response
func ErrorRes(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(IMap{
		"error": err.Error(),
	})
}
