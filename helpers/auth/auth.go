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

// TODO we should avoid string maps as they are slow
type Auth struct {
	keys     map[string]key
	baseSeed []byte
}

type key struct {
	keyBytes []byte
	apiKey   models.ApiKey
	sha512   []rollingHash
	sha256   []rollingHash
}

type rollingHash struct {
	salt  []byte
	value []byte
}

func New(keys []models.ApiKey, baseSeed []byte) *Auth {
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

var ErrorInvalidKey = errors.New("invalid key")

func (a *Auth) GetBaseSeed() []byte {
	return a.baseSeed
}

func (a *Auth) Authenticate(authorizationHeader []byte) (site *models.ApiKey, salt []byte, err error) {
	authorizationHeaderLen := len(authorizationHeader)
	if authorizationHeaderLen < 7 {
		if authorizationHeaderLen == 0 {
			return nil, salt, errors.New("missing authorization header of type Basic")
		}
		return nil, salt, errors.New("invalid authorization header, must be of type Basic and contain data")
	}

	if !bytes.Equal([]byte("Basic "), authorizationHeader[:6]) {
		return nil, salt, errors.New("authorization must be of type Basic")
	}

	auth, err := base64.RawStdEncoding.DecodeString(string(authorizationHeader[6:]))
	if err != nil {
		return nil, salt, err
	}

	parts := bytes.Split(auth, []byte(":"))
	if len(parts) != 4 {
		return nil, salt, ErrorInvalidKey
	}

	isSha512 := bytes.Equal(parts[0], []byte("sha512"))
	if !isSha512 && !bytes.Equal(parts[0], []byte("sha256")) {
		return nil, salt, errors.New("only sha512 and sha256 are supported")
	}

	siteId := string(parts[1])
	if !primitive.IsValidObjectID(siteId) {
		return nil, salt, errors.New("invalid site ID")
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

	knownKey, ok := a.keys[siteId]
	if !ok {
		return nil, salt, ErrorInvalidKey
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
				return nil, salt, ErrorInvalidKey
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
			return nil, salt, ErrorInvalidKey
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
				return nil, salt, ErrorInvalidKey
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
			return nil, salt, ErrorInvalidKey
		}

		// Pre calculate next hash in the chain
		hash = sha256.Sum256(append(hash[:], keyAndSalt...))

		knownKey.sha256 = append(knownKey.sha256, rollingHash{
			salt:  salt,
			value: hash[:],
		})
	}

	a.keys[siteId] = knownKey
	return &knownKey.apiKey, salt, nil
}
