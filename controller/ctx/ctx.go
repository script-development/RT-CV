package ctx

import (
	"context"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

type profilesCtx uint8
type profileCtx uint8
type authCtx uint8
type keyCtx uint8
type loggerCtx uint8
type dbConnCtx uint8

const (
	profilesCtxKey = profilesCtx(0)
	profileCtxKey  = profileCtx(0)
	authCtxKey     = authCtx(0)
	keyCtxKey      = keyCtx(0)
	loggerCtxKey   = loggerCtx(0)
	dbConnCtxKey   = dbConnCtx(0)
)

func getCtxValue(c *fiber.Ctx, key interface{}) interface{} {
	return c.UserContext().Value(key)
}

// GetProfiles returns the search profiles
func GetProfiles(c *fiber.Ctx) *[]models.Profile {
	return getCtxValue(c, profilesCtxKey).(*[]models.Profile)
}

// SetProfiles sets the search profiles
func SetProfiles(ctx context.Context, value *[]models.Profile) context.Context {
	return context.WithValue(ctx, profilesCtxKey, value)
}

// GetProfile returns a single profile
func GetProfile(c *fiber.Ctx) *models.Profile {
	return getCtxValue(c, profileCtxKey).(*models.Profile)
}

// SetProfile sets a single search profile
func SetProfile(ctx context.Context, value *models.Profile) context.Context {
	return context.WithValue(ctx, profileCtxKey, value)
}

// GetAuth returns the auth key used to make the request
// If the request is not authenticated this function panics
func GetAuth(c *fiber.Ctx) *auth.Auth {
	return getCtxValue(c, authCtxKey).(*auth.Auth)
}

// SetAuth sets the auth key used to make the request
func SetAuth(ctx context.Context, value *auth.Auth) context.Context {
	return context.WithValue(ctx, authCtxKey, value)
}

// GetKey returns the api key used to make the request
func GetKey(c *fiber.Ctx) *models.APIKey {
	return getCtxValue(c, keyCtxKey).(*models.APIKey)
}

// SetKey sets the api key used to make the request
func SetKey(ctx context.Context, value *models.APIKey) context.Context {
	return context.WithValue(ctx, keyCtxKey, value)
}

// GetLogger returns the global logger
func GetLogger(c *fiber.Ctx) *log.Entry {
	return getCtxValue(c, loggerCtxKey).(*log.Entry)
}

// SetLogger sets the logger used in the request
func SetLogger(ctx context.Context, value *log.Entry) context.Context {
	return context.WithValue(ctx, loggerCtxKey, value)
}

// GetDbConn returns the database connection
func GetDbConn(c *fiber.Ctx) db.Connection {
	return getCtxValue(c, dbConnCtxKey).(db.Connection)
}

// SetDbConn sets the database connection
func SetDbConn(ctx context.Context, value db.Connection) context.Context {
	return context.WithValue(ctx, dbConnCtxKey, value)
}
