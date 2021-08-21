package ctx

import (
	"context"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

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

func getCtxValue(c *fiber.Ctx, key interface{}) interface{} {
	return c.UserContext().Value(key)
}

func GetProfiles(c *fiber.Ctx) *[]models.Profile {
	return getCtxValue(c, profilesCtxKey).(*[]models.Profile)
}
func SetProfiles(ctx context.Context, value *[]models.Profile) context.Context {
	return context.WithValue(ctx, profilesCtxKey, value)
}

func GetAuth(c *fiber.Ctx) *auth.Auth {
	return getCtxValue(c, authCtxKey).(*auth.Auth)
}
func SetAuth(ctx context.Context, value *auth.Auth) context.Context {
	return context.WithValue(ctx, profilesCtxKey, value)
}

func GetKey(c *fiber.Ctx) *models.ApiKey {
	return getCtxValue(c, keyCtxKey).(*models.ApiKey)
}
func SetKey(ctx context.Context, value *models.ApiKey) context.Context {
	return context.WithValue(ctx, profilesCtxKey, value)
}

func GetLogger(c *fiber.Ctx) *log.Entry {
	return getCtxValue(c, loggerCtxKey).(*log.Entry)
}
func SetLogger(ctx context.Context, value *log.Entry) context.Context {
	return context.WithValue(ctx, profilesCtxKey, value)
}

func GetDbConn(c *fiber.Ctx) dbInterfaces.Connection {
	return getCtxValue(c, dbConnCtxKey).(dbInterfaces.Connection)
}
func SetDbConn(ctx context.Context, value dbInterfaces.Connection) context.Context {
	return context.WithValue(ctx, profilesCtxKey, value)
}
