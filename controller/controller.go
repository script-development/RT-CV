package controller

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

func Routes(app *fiber.App, dbConn db.Connection) {
	v1 := app.Group(`/v1`, InsertData(dbConn))

	authRoutes(v1)
	scraperRoutes(v1)
	controllerRoutes(v1)
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
func GetDbConn(c *fiber.Ctx) db.Connection {
	return GetCtxValue(c, dbConnCtxKey).(db.Connection)
}
