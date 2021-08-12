package auth

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/script-development/RT-CV/models"
)

type Auth map[int]key

type key struct {
	siteId int
	key    []byte
	sha512 []rollingHash
	sha256 []rollingHash
}

type rollingHash struct {
	salt  []byte
	value []byte
}

func New(keys []models.ApiKey) *Auth {
	res := Auth{}
	for _, dbKey := range keys {
		if !dbKey.Enabled {
			continue
		}
		res[dbKey.ID] = key{
			siteId: dbKey.SiteId,
			key:    []byte(dbKey.Key),
			sha512: []rollingHash{},
			sha256: []rollingHash{},
		}
	}
	return &res
}

var errorInvalidKey = errors.New("invalid key")

func (a *Auth) Authenticate(authorizationHeader []byte) (siteID int, err error) {
	if len(authorizationHeader) < 7 {
		return 0, errors.New("invalid authorization value, must be of type Basic and contain data")
	}

	if !bytes.Equal([]byte("Basic "), authorizationHeader[:6]) {
		return 0, errors.New("authorization must be of type Basic")
	}

	auth, err := base64.RawStdEncoding.DecodeString(string(authorizationHeader[6:]))
	if err != nil {
		return 0, err
	}

	parts := bytes.Split(auth, []byte(":"))
	if len(parts) != 4 {
		return 0, errorInvalidKey
	}

	isSha512 := bytes.Equal(parts[0], []byte("sha512"))
	if !isSha512 && !bytes.Equal(parts[0], []byte("sha256")) {
		return 0, errors.New("only sha512 and sha256 are supported")
	}

	siteId, err := strconv.Atoi(string(parts[1]))
	if err != nil || siteId < 1 {
		return 0, errors.New("site id is not a positive number")
	}

	salt := parts[2]
	if len(salt) == 0 {
		return 0, errors.New("salt cannot be empty")
	}

	key := parts[3]
	if len(key) == 0 {
		return 0, errors.New("key cannot be empty")
	}
	n, err := hex.Decode(key, key)
	if err != nil {
		return 0, errors.New("invalid key hash")
	}
	key = key[:n]

	knownKey, ok := (*a)[siteId]
	if !ok {
		return 0, errorInvalidKey
	}

	if isSha512 {
		for idx, entry := range knownKey.sha512 {
			if !bytes.Equal(entry.salt, salt) {
				continue
			}
			if !bytes.Equal(entry.value, key) {
				return 0, errorInvalidKey
			}

			hash := sha512.Sum512(append(append(entry.value, knownKey.key...), entry.salt...))
			entry.value = hash[:]
			knownKey.sha512[idx] = entry

			return knownKey.siteId, nil
		}

		hash := sha512.Sum512(append(knownKey.key, salt...))
		hash = sha512.Sum512(append(append(hash[:], knownKey.key...), salt...))

		if !bytes.Equal(hash[:], key) {
			return 0, errorInvalidKey
		}

		// Pre calculate next hash in the chain
		hash = sha512.Sum512(append(append(hash[:], knownKey.key...), salt...))

		knownKey.sha512 = append(knownKey.sha512, rollingHash{
			salt:  salt,
			value: hash[:],
		})
	} else {
		for idx, entry := range knownKey.sha256 {
			if !bytes.Equal(entry.salt, salt) {
				continue
			}
			if !bytes.Equal(entry.value, key) {
				return 0, errorInvalidKey
			}

			hash := sha256.Sum256(append(append(entry.value, knownKey.key...), entry.salt...))
			entry.value = hash[:]
			knownKey.sha256[idx] = entry

			return knownKey.siteId, nil
		}

		hash := sha256.Sum256(append(knownKey.key, salt...))
		hash = sha256.Sum256(append(append(hash[:], knownKey.key...), salt...))

		if !bytes.Equal(hash[:], key) {
			return 0, errorInvalidKey
		}

		// Pre calculate next hash in the chain
		hash = sha256.Sum256(append(append(hash[:], knownKey.key...), salt...))

		knownKey.sha256 = append(knownKey.sha256, rollingHash{
			salt:  salt,
			value: hash[:],
		})
	}

	(*a)[siteId] = knownKey
	return knownKey.siteId, nil
}
