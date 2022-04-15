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

type ctx uint8

const ctxKey = ctx(0)

// Ctx contains the request context
type Ctx struct {
	RequestID            primitive.ObjectID
	Profile              *models.Profile
	Auth                 *auth.Helper
	Key                  *models.APIKey
	APIKeyFromParam      *models.APIKey
	Logger               *log.Entry
	DBConn               db.Connection
	MatcherProfilesCache *MatcherProfilesCache
	OnMatchHook          *models.OnMatchHook
}

// Set sets the request context
func Set(ctx context.Context, value *Ctx) context.Context {
	return context.WithValue(ctx, ctxKey, value)
}

// Get returns the request context
func Get(c *fiber.Ctx) *Ctx {
	return c.UserContext().Value(ctxKey).(*Ctx)
}

// MatcherProfilesCache contains the matcher profiles cache
type MatcherProfilesCache struct {
	InsertionTime time.Time
	Profiles      []*models.Profile
}

// ResetMatcherProfilesCache sets the profiles cache to an empty object
func (c *Ctx) ResetMatcherProfilesCache() {
	*c.MatcherProfilesCache = MatcherProfilesCache{}
}
