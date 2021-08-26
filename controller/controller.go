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
		Group(c, `/schema`, func(c fiber.Router) {
			c.Get(`/cv`, routeGetCvSchema)
		})

		Group(c, `/auth`, func(c fiber.Router) {
			c.Get(`/keyinfo`, requiresAuth(0), routeGetKeyInfo)
			c.Get(`/seed`, routeAuthSeed)
		})

		Group(c, `/scraper`, func(c fiber.Router) {
			c.Post(`/scanCV`, routeScraperScanCV)
		}, requiresAuth(models.APIKeyRoleScraper|models.APIKeyRoleDashboard))

		secretsRoutes := func(c fiber.Router) {
			c.Get(``, routeGetSecrets)
			Group(c, `/:key`, func(c fiber.Router) {
				c.Delete(``, routeDeleteSecret)
				Group(c, `/:encryptionKey`, func(c fiber.Router) {
					c.Get(``, routeGetSecret)
					c.Put(``, routeUpdateSecret)
					c.Post(``, routeCreateSecret)
				}, validEncryptionKeyMiddleware())
			}, validKeyMiddleware())
		}
		Group(c, `/secrets/myKey`, secretsRoutes, requiresAuth(models.APIKeyRoleAll), middlewareBindMyKey())
		Group(c, `/secrets/otherKey`, func(c fiber.Router) {
			c.Get(``, routeGetAllSecretsFromAllKeys)
			Group(c, `/:keyID`, secretsRoutes, middlewareBindKey())
		}, requiresAuth(models.APIKeyRoleDashboard))

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
			c.Post(``, routeCreateKey)
			Group(c, `/:keyID`, func(c fiber.Router) {
				c.Get(``, routeGetKey)
				c.Put(``, routeUpdateKey)
				c.Put(``, routeUpdateKey)
				c.Delete(``, routeDeleteKey)
			}, middlewareBindKey())
		}, requiresAuth(models.APIKeyRoleDashboard))
	}, InsertData(dbConn, serverSeed))

	app.Static("", "./dashboard/out/index.html", fiber.Static{Compress: true})
	app.Static("login", "./dashboard/out/login.html", fiber.Static{Compress: true})
	app.Static("_next", "./dashboard/out/_next", fiber.Static{Compress: true})
	app.Static("favicon.ico", "./dashboard/out/favicon.ico", fiber.Static{Compress: true})
	app.Use(func(c *fiber.Ctx) error {
		// 404 page
		return c.Status(404).SendFile("./dashboard/out/404.html", true)
	})
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
