package auth

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hashMethod uint8

var (
	hashSha512 = hashMethod(0)
	hashSha256 = hashMethod(1)
)

func hashSha(method hashMethod, in string) string {
	if method == hashSha256 {
		res := sha256.Sum256([]byte(in))
		return hex.EncodeToString(res[:])
	}

	res := sha512.Sum512([]byte(in))
	return hex.EncodeToString(res[:])
}

// Auth can be used to check authentication headers
type Auth struct {
	keys     map[string]key
	baseSeed string
}

type key struct {
	key    string
	apiKey models.APIKey
	sha512 []rollingHash
	sha256 []rollingHash
}

type rollingHash struct {
	salt  string
	value string
}

// New returns a new Auth instance that can be used to check auth tokens
func New(keys []models.APIKey, baseSeed string) *Auth {
	res := Auth{
		keys:     map[string]key{},
		baseSeed: baseSeed,
	}
	for _, dbKey := range keys {
		if !dbKey.Enabled {
			continue
		}

		res.keys[dbKey.ID.Hex()] = key{
			apiKey: dbKey,
			key:    dbKey.Key,
			sha512: []rollingHash{},
			sha256: []rollingHash{},
		}
	}
	return &res
}

var (
	// ErrInvalidKey is a returned when a key is invalid
	ErrInvalidKey = errors.New("invalid authentication key")
	// ErrNoAuthheader is send when the authentication header is empty
	ErrNoAuthheader = errors.New("missing authorization header of type Basic")
	// ErrAuthHeaderToShort is send whenthe authorization is to short to even check
	ErrAuthHeaderToShort = errors.New("invalid authorization header, must be of type Basic and contain data")
)

// GetBaseSeed returns the server base seed
func (a *Auth) GetBaseSeed() string {
	return a.baseSeed
}

// RefreshKey resets a key
// This should be called if the updatedKey.key changed otherwise users can still send auth requests with the old key
func (a *Auth) RefreshKey(updatedKey models.APIKey) {
	id := updatedKey.ID.Hex()
	if !updatedKey.Enabled {
		delete(a.keys, id)
		return
	}

	a.keys[id] = key{
		apiKey: updatedKey,
		key:    updatedKey.Key,
		sha512: []rollingHash{},
		sha256: []rollingHash{},
	}
}

// AddKey Adds a key to the authenticator so it can be used to authenticate with
func (a *Auth) AddKey(newKey models.APIKey) {
	if !newKey.Enabled {
		return
	}

	a.keys[newKey.ID.Hex()] = key{
		apiKey: newKey,
		key:    newKey.Key,
		sha512: []rollingHash{},
		sha256: []rollingHash{},
	}
}

// Authenticate check is a authorizationHeader is correct
func (a *Auth) Authenticate(authorizationHeader string) (site *models.APIKey, salt string, err error) {
	authorizationHeaderLen := len(authorizationHeader)
	if authorizationHeaderLen < 7 {
		if authorizationHeaderLen == 0 {
			return nil, salt, ErrNoAuthheader
		}
		return nil, salt, ErrAuthHeaderToShort
	}

	if "Basic " != authorizationHeader[:6] {
		return nil, salt, errors.New("authorization must be of type Basic")
	}

	auth, err := base64.URLEncoding.DecodeString(authorizationHeader[6:])
	if err != nil {
		return nil, salt, err
	}

	parts := strings.Split(string(auth), ":")
	if len(parts) != 4 {
		return nil, salt, ErrInvalidKey
	}

	hashMethod := hashSha256
	if parts[0] == "sha512" {
		hashMethod = hashSha512
	}
	if hashMethod == hashSha256 && parts[0] != "sha256" {
		return nil, salt, errors.New("only sha512 and sha256 are supported")
	}

	siteID := string(parts[1])
	if !primitive.IsValidObjectID(siteID) {
		return nil, salt, errors.New("invalid key ID")
	}

	salt = parts[2]
	if len(salt) == 0 {
		return nil, salt, errors.New("salt cannot be empty")
	}

	key := parts[3]
	if len(key) == 0 {
		return nil, salt, errors.New("key cannot be empty")
	}

	knownKey, ok := a.keys[siteID]
	if !ok {
		return nil, salt, ErrInvalidKey
	}
	keyAndSalt := knownKey.key + salt

	itemsArr := knownKey.sha512
	if hashMethod == hashSha256 {
		itemsArr = knownKey.sha256
	}

	for idx, entry := range itemsArr {
		if entry.salt != salt {
			continue
		}

		// Key + salt combo earlier created lets check if the credentials match
		if entry.value != key {
			return nil, salt, ErrInvalidKey
		}

		entry.value = hashSha(hashMethod, entry.value+knownKey.key+entry.salt)
		if hashMethod == hashSha512 {
			knownKey.sha512[idx] = entry
		} else {
			knownKey.sha256[idx] = entry
		}

		return &knownKey.apiKey, salt, nil
	}

	// Create a new key + salt combo
	hash := hashSha(hashMethod, a.baseSeed+keyAndSalt)
	hash = hashSha(hashMethod, hash+keyAndSalt)

	if hash != key {
		return nil, salt, ErrInvalidKey
	}

	// Pre calculate next hash in the chain
	hash = hashSha(hashMethod, hash+keyAndSalt)

	rollingKey := rollingHash{
		salt:  salt,
		value: hash,
	}
	if hashMethod == hashSha512 {
		knownKey.sha512 = append(knownKey.sha512, rollingKey)
	} else {
		knownKey.sha256 = append(knownKey.sha256, rollingKey)
	}

	a.keys[siteID] = knownKey
	return &knownKey.apiKey, salt, nil
}
