package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/puzpuzpuz/xsync"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/crypto"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GenAuthHeaderKey generates a authentication header value
func GenAuthHeaderKey(id, key string) string {
	return "Basic " + id + ":" + key
}

// Helper helps authenticate a user
type Helper struct {
	// the cache key is the key ID
	cache  *xsync.MapOf[string, cachedKey]
	dbConn db.Connection
}

type cachedKey struct {
	validTil                                        time.Time
	KeyAsString, KeyAsSha512Lower, KeyAsSha512Upper string
	key                                             *models.APIKey
}

// NewHelper returns a new instance of AuthHelper
func NewHelper(dbConn db.Connection) *Helper {
	return &Helper{
		cache:  xsync.NewMapOf[cachedKey](),
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
	var keyIsSha512Hashed bool
	switch len(authorizationHeader) {
	case 159:
		// header with Sha512 key
		keyIsSha512Hashed = true
	case 63:
		// header with raw key
		keyIsSha512Hashed = false
	default:
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

	authKey := authorizationHeader[endID+1:]
	keyCacheEntry, ok := h.cache.Load(id)
	if ok {
		if time.Now().Before(keyCacheEntry.validTil) {
			// Yay a cache entry for this key exists and it's still valid
			if keyIsSha512Hashed {
				if keyCacheEntry.KeyAsSha512Lower == authKey || keyCacheEntry.KeyAsSha512Upper == authKey {
					return keyCacheEntry.key, nil
				}
			} else {
				if keyCacheEntry.KeyAsString == authKey {
					return keyCacheEntry.key, nil
				}
			}
			return nil, ErrAuthHeaderInvalid
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

	hashedKey := strings.ToLower(crypto.HashSha512String(key.Key))
	hashedKeyUpper := strings.ToUpper(hashedKey)

	h.cache.Store(id, cachedKey{
		validTil:         time.Now().Add(time.Hour * 12),
		KeyAsString:      key.Key,
		KeyAsSha512Lower: strings.ToLower(hashedKey),
		KeyAsSha512Upper: hashedKeyUpper,
		key:              &key,
	})

	if keyIsSha512Hashed {
		if authKey == hashedKey || authKey == hashedKeyUpper {
			return &key, nil
		}
	} else {
		if authKey == key.Key {
			return &key, nil
		}
	}
	return nil, ErrAuthHeaderInvalid
}
