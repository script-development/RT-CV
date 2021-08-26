package auth

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth can be used to check authentication headers
type Auth struct {
	keys     map[string]key
	baseSeed []byte
}

type key struct {
	keyBytes []byte
	apiKey   models.APIKey
	sha512   []rollingHash
	sha256   []rollingHash
}

type rollingHash struct {
	salt  []byte
	value []byte
}

// New returns a new Auth instance that can be used to check auth tokens
func New(keys []models.APIKey, baseSeed []byte) *Auth {
	res := Auth{
		keys:     map[string]key{},
		baseSeed: baseSeed,
	}
	for _, dbKey := range keys {
		if !dbKey.Enabled {
			continue
		}

		res.keys[dbKey.ID.Hex()] = key{
			apiKey:   dbKey,
			keyBytes: []byte(dbKey.Key),
			sha512:   []rollingHash{},
			sha256:   []rollingHash{},
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
func (a *Auth) GetBaseSeed() []byte {
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
		apiKey:   updatedKey,
		keyBytes: []byte(updatedKey.Key),
		sha512:   []rollingHash{},
		sha256:   []rollingHash{},
	}
}

// AddKey Adds a key to the authenticator so it can be used to authenticate with
func (a *Auth) AddKey(newKey models.APIKey) {
	if !newKey.Enabled {
		return
	}

	a.keys[newKey.ID.Hex()] = key{
		apiKey:   newKey,
		keyBytes: []byte(newKey.Key),
		sha512:   []rollingHash{},
		sha256:   []rollingHash{},
	}
}

// Authenticate check is a authorizationHeader is correct
func (a *Auth) Authenticate(authorizationHeader []byte) (site *models.APIKey, salt []byte, err error) {
	authorizationHeaderLen := len(authorizationHeader)
	if authorizationHeaderLen < 7 {
		if authorizationHeaderLen == 0 {
			return nil, salt, ErrNoAuthheader
		}
		return nil, salt, ErrAuthHeaderToShort
	}

	if !bytes.Equal([]byte("Basic "), authorizationHeader[:6]) {
		return nil, salt, errors.New("authorization must be of type Basic")
	}

	auth, err := base64.URLEncoding.DecodeString(string(authorizationHeader[6:]))
	if err != nil {
		return nil, salt, err
	}

	parts := bytes.Split(auth, []byte(":"))
	if len(parts) != 4 {
		return nil, salt, ErrInvalidKey
	}

	isSha512 := bytes.Equal(parts[0], []byte("sha512"))
	if !isSha512 && !bytes.Equal(parts[0], []byte("sha256")) {
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
	n, err := hex.Decode(key, key)
	if err != nil {
		return nil, salt, errors.New("invalid key hash")
	}
	key = key[:n]

	knownKey, ok := a.keys[siteID]
	if !ok {
		return nil, salt, ErrInvalidKey
	}
	keyAndSalt := append(knownKey.keyBytes, salt...)

	// FIXME: The sha512 code is almost equal to that of sha256
	if isSha512 {
		for idx, entry := range knownKey.sha512 {
			if !bytes.Equal(entry.salt, salt) {
				continue
			}

			// Key + salt combo earlier created lets check if the credentials match
			if !bytes.Equal(entry.value, key) {
				return nil, salt, ErrInvalidKey
			}

			hash := sha512.Sum512(append(append(entry.value, knownKey.keyBytes...), entry.salt...))
			entry.value = hash[:]
			knownKey.sha512[idx] = entry

			return &knownKey.apiKey, salt, nil
		}

		// Create a new key + salt combo
		hash := sha512.Sum512(append(a.baseSeed, keyAndSalt...))
		hash = sha512.Sum512(append(hash[:], keyAndSalt...))

		if !bytes.Equal(hash[:], key) {
			return nil, salt, ErrInvalidKey
		}

		// Pre calculate next hash in the chain
		hash = sha512.Sum512(append(hash[:], keyAndSalt...))

		knownKey.sha512 = append(knownKey.sha512, rollingHash{
			salt:  salt,
			value: hash[:],
		})
	} else {
		for idx, entry := range knownKey.sha256 {
			if !bytes.Equal(entry.salt, salt) {
				continue
			}

			// Key + salt combo earlier created lets check if the credentials match
			if !bytes.Equal(entry.value, key) {
				return nil, salt, ErrInvalidKey
			}

			hash := sha256.Sum256(append(entry.value, keyAndSalt...))
			entry.value = hash[:]
			knownKey.sha256[idx] = entry

			return &knownKey.apiKey, salt, nil
		}

		// Create a new key + salt combo
		hash := sha256.Sum256(append(a.baseSeed, keyAndSalt...))
		hash = sha256.Sum256(append(hash[:], keyAndSalt...))

		if !bytes.Equal(hash[:], key) {
			return nil, salt, ErrInvalidKey
		}

		// Pre calculate next hash in the chain
		hash = sha256.Sum256(append(hash[:], keyAndSalt...))

		knownKey.sha256 = append(knownKey.sha256, rollingHash{
			salt:  salt,
			value: hash[:],
		})
	}

	a.keys[siteID] = knownKey
	return &knownKey.apiKey, salt, nil
}
