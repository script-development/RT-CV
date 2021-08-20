package controller

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

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
				c.Delete(``, routeScraperDeleteSecret)
				Group(c, `/:secretKey`, func(c fiber.Router) {
					c.Post(``, routeScraperCreateSecret)
					c.Put(``, routeScraperUpdateSecret)
					c.Get(``, routeScraperGetSecret)
				}, validSecretKey())
			}, validKey())
		}, requiresAuth(models.ApiKeyRoleScraper))

		Group(c, `/control`, func(c fiber.Router) {
			c.Get(`/reloadProfiles`, routeControlReloadProfiles)
		}, requiresAuth(models.ApiKeyRoleAdmin|models.ApiKeyRoleController))
	}, InsertData(dbConn, serverSeed))
}

type profilesCtx uint8
type authCtx uint8
type keyCtx uint8
type loggerCtx uint8
type dbConnCtx uint8

const (
	profilesCtxKey = profilesCtx(0)
	authCtxKey     = authCtx(0)
	keyCtxKey      = keyCtx(0)
	loggerCtxKey   = loggerCtx(0)
	dbConnCtxKey   = dbConnCtx(0)
)

func GetCtxValue(c *fiber.Ctx, key interface{}) interface{} {
	return c.UserContext().Value(key)
}
func GetProfiles(c *fiber.Ctx) *[]models.Profile {
	return GetCtxValue(c, profilesCtxKey).(*[]models.Profile)
}
func GetAuth(c *fiber.Ctx) *auth.Auth {
	return GetCtxValue(c, authCtxKey).(*auth.Auth)
}
func GetKey(c *fiber.Ctx) *models.ApiKey {
	return GetCtxValue(c, keyCtxKey).(*models.ApiKey)
}
func GetLogger(c *fiber.Ctx) *log.Entry {
	return GetCtxValue(c, loggerCtxKey).(*log.Entry)
}
func GetDbConn(c *fiber.Ctx) dbInterfaces.Connection {
	return GetCtxValue(c, dbConnCtxKey).(dbInterfaces.Connection)
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	return ErrorRes(c, 500, err)
}

func ErrorRes(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(map[string]string{
		"error": err.Error(),
	})
}
