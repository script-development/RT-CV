package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/crypto"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GenAuthHeaderKey generates a authentication header value
func GenAuthHeaderKey(id, key string) string {
	return "Basic " + id + ":" + crypto.HashSha512String(key)
}

// Helper helps authenticate a user
type Helper struct {
	// the cache key is the key ID
	cache  sync.Map // = map[string]cachedKey
	dbConn db.Connection
}

type cachedKey struct {
	LastRefreshed time.Time
	KeyAsSha512   string
	key           models.APIKey
}

// NewHelper returns a new instance of AuthHelper
func NewHelper(dbConn db.Connection) *Helper {
	return &Helper{
		dbConn: dbConn,
	}
}

var (
	// ErrNoAuthHeader = no authorization header
	ErrNoAuthHeader = errors.New("missing authorization header of type Basic")
	// ErrAuthHeaderHasInvalidLen = auth header has invalid length
	ErrAuthHeaderHasInvalidLen = errors.New("auth header has invalid length")
	// ErrAuthHeaderNotBasic = auth header expected to be \"Basic ...\"
	ErrAuthHeaderNotBasic = errors.New("auth header expected to be \"Basic ...\"")
	// ErrAuthHeaderInvalidFormat = auth header has invalid format, expect \"Basic keyID:sha512(Key)\"
	ErrAuthHeaderInvalidFormat = errors.New("auth header has invalid format, expect \"Basic keyID:sha512(Key)\"")
	// ErrAuthHeaderInvalid = auth header is invalid
	ErrAuthHeaderInvalid = errors.New("auth header is invalid")
)

// RemoveKeyCache removes a cached key
func (h *Helper) RemoveKeyCache(id string) {
	h.cache.Delete(id)
}

// Valid validates an authorizationHeader
func (h *Helper) Valid(authorizationHeader string) (*models.APIKey, error) {
	if len(authorizationHeader) != 159 {
		return nil, ErrAuthHeaderHasInvalidLen
	}

	if "Basic " != authorizationHeader[:6] {
		return nil, ErrAuthHeaderNotBasic
	}

	startID := 6
	endID := 6 + 24

	id := authorizationHeader[startID:endID]

	if authorizationHeader[endID] != ':' {
		return nil, ErrAuthHeaderInvalidFormat
	}

	keyAsSha512 := authorizationHeader[endID+1:]

	keyCacheEntryInterf, ok := h.cache.Load(id)
	if ok {
		keyCacheEntry := keyCacheEntryInterf.(cachedKey)
		if time.Now().Before(keyCacheEntry.LastRefreshed.Add(time.Hour * 12)) {
			// Yay a cache entry for this key exists and is still valid
			if keyCacheEntry.KeyAsSha512 != keyAsSha512 {
				return nil, ErrAuthHeaderInvalid
			}
			return &keyCacheEntry.key, nil
		}
		// Cache entry outdated
		h.RemoveKeyCache(id)
	}

	parsedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrAuthHeaderInvalidFormat
	}
	key, err := models.GetAPIKey(h.dbConn, parsedID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrAuthHeaderInvalid
		}
		return nil, err
	}

	hashedKey := crypto.HashSha512String(key.Key)
	h.cache.Store(id, cachedKey{
		LastRefreshed: time.Now(),
		KeyAsSha512:   hashedKey,
		key:           key,
	})

	if keyAsSha512 != hashedKey {
		return nil, ErrAuthHeaderInvalid
	}
	return &key, nil
}
