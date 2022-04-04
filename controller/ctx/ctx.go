package ctx

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
// Maybe we should should replace the below for only one value as in 1 struct with all the items below
// It might give a performance improvement and also make the code more safe as then we can do nil checks
//

type profileCtx uint8
type authCtx uint8
type keyCtx uint8
type keyFromParamCtx uint8
type loggerCtx uint8
type dbConnCtx uint8
type requestIDCtx uint8
type profilesCacheCtx uint8
type onMatchHookKey uint8

const (
	profileCtxKey       = profileCtx(0)
	authCtxKey          = authCtx(0)
	keyCtxKey           = keyCtx(0)
	keyFromParamCtxKey  = keyFromParamCtx(0)
	loggerCtxKey        = loggerCtx(0)
	dbConnCtxKey        = dbConnCtx(0)
	requestIDCtxKey     = requestIDCtx(0)
	profilesCacheCtxKey = profilesCacheCtx(0)
	onMatchHookCtxKey   = onMatchHookKey(0)
)

// getCtxValue returns a value from the context
// Will panic if the key is not found
func getCtxValue(c *fiber.Ctx, key any) any {
	return c.UserContext().Value(key)
}

// GetRequestID returns the request id
// We bind this to every request so we can debug things more easily in the case of many requests
func GetRequestID(c *fiber.Ctx) primitive.ObjectID {
	return getCtxValue(c, requestIDCtxKey).(primitive.ObjectID)
}

// SetRequestID sets the request id
func SetRequestID(ctx context.Context, value primitive.ObjectID) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, value)
}

// GetProfile returns a profile (probably based on the profile id specified in the url)
func GetProfile(c *fiber.Ctx) *models.Profile {
	return getCtxValue(c, profileCtxKey).(*models.Profile)
}

// SetProfile sets a profile
func SetProfile(ctx context.Context, value *models.Profile) context.Context {
	return context.WithValue(ctx, profileCtxKey, value)
}

// GetAuth returns the auth key used to make the request
// If the request is not authenticated this function panics
func GetAuth(c *fiber.Ctx) *auth.Helper {
	return getCtxValue(c, authCtxKey).(*auth.Helper)
}

// SetAuth sets the auth key used to make the request
func SetAuth(ctx context.Context, value *auth.Helper) context.Context {
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

// GetAPIKeyFromParam returns the api key specified in the url
func GetAPIKeyFromParam(c *fiber.Ctx) *models.APIKey {
	return getCtxValue(c, keyFromParamCtxKey).(*models.APIKey)
}

// SetAPIKeyFromParam sets an api key based on the route
func SetAPIKeyFromParam(ctx context.Context, value *models.APIKey) context.Context {
	return context.WithValue(ctx, keyFromParamCtxKey, value)
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

// MatcherProfilesCache contains the matcher profiles cache
type MatcherProfilesCache struct {
	InsertionTime time.Time
	Profiles      []*models.Profile
}

// GetMatcherProfilesCache returns the cached profiles used by the matcher
func GetMatcherProfilesCache(c *fiber.Ctx) *MatcherProfilesCache {
	return getCtxValue(c, profilesCacheCtxKey).(*MatcherProfilesCache)
}

// ResetMatcherProfilesCache sets the profiles cache to an empty object
func ResetMatcherProfilesCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, profilesCacheCtxKey, &MatcherProfilesCache{})
}

// SetOnMatchHook sets the on match hook in the request context
func SetOnMatchHook(ctx context.Context, value *models.OnMatchHook) context.Context {
	return context.WithValue(ctx, onMatchHookCtxKey, value)
}

// GetOnMatchHook returns the on match hook from the request context
func GetOnMatchHook(c *fiber.Ctx) *models.OnMatchHook {
	return getCtxValue(c, onMatchHookCtxKey).(*models.OnMatchHook)
}
