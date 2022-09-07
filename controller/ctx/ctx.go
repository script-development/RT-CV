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
	Key                  *models.APIKey // The key used to make the request
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

// GetOrGenMatcherProfilesCache returns the cached profiles or creates a new cache content
func (c *Ctx) GetOrGenMatcherProfilesCache() (*MatcherProfilesCache, error) {
	if c.MatcherProfilesCache.ScanProfiles != nil && c.MatcherProfilesCache.ListProfiles != nil && c.MatcherProfilesCache.InsertionTime.Add(time.Hour*24).After(time.Now()) {
		return c.MatcherProfilesCache, nil
	}

	// Update the cache
	c.Logger.Info("updating the profiles cache")
	scanProfiles, err := models.GetActualActiveProfiles(c.DBConn)
	if err != nil {
		return nil, err
	}
	listProfiles, err := models.GetListsProfiles(c.DBConn)
	if err != nil {
		return nil, err
	}

	*c.MatcherProfilesCache = MatcherProfilesCache{
		ScanProfiles:  profilesListToPtrs(scanProfiles),
		ListProfiles:  profilesListToPtrs(listProfiles),
		InsertionTime: time.Now(),
	}
	return c.MatcherProfilesCache, nil
}

func profilesListToPtrs(in []models.Profile) []*models.Profile {
	out := make([]*models.Profile, len(in))
	for idx := range in {
		out[idx] = &in[idx]
	}
	return out
}

// MatcherProfilesCache contains the matcher profiles cache
type MatcherProfilesCache struct {
	InsertionTime time.Time
	ScanProfiles  []*models.Profile
	ListProfiles  []*models.Profile
}

// ResetMatcherProfilesCache sets the profiles cache to an empty object
func (c *Ctx) ResetMatcherProfilesCache() {
	*c.MatcherProfilesCache = MatcherProfilesCache{}
}
