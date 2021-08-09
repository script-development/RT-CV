package auth

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
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

type Keys struct{}

type DbKey struct {
	ID     int
	SiteID int
	Key    string
}

func New(keys []DbKey) *Auth {
	res := Auth{}
	for _, dbKey := range keys {
		res[dbKey.ID] = key{
			key:    []byte(dbKey.Key),
			sha512: []rollingHash{},
			sha256: []rollingHash{},
		}
	}
	return &res
}

var errorInvalidKey = errors.New("invalid key")

func (a *Auth) Authenticate(authorizationHeader []byte) (siteID int, err error) {
	if !bytes.Equal([]byte("Basic "), authorizationHeader) {
		return 0, errors.New("authorization must be of type basic")
	}

	auth := make([]byte, base64.RawStdEncoding.DecodedLen(len(authorizationHeader)-6))
	n, err := base64.RawStdEncoding.Decode(auth, authorizationHeader[6:])
	if err != nil {
		return 0, err
	}
	auth = auth[:n]

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
	n, err = hex.Decode(key, key)
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

			hash := sha512.Sum512(entry.value)
			entry.value = hash[:]
			knownKey.sha512[idx] = entry

			return knownKey.siteId, nil
		}

		hash := sha512.Sum512(append(knownKey.key, salt...))
		hash = sha512.Sum512(hash[:])

		if !bytes.Equal(hash[:], key) {
			return 0, errorInvalidKey
		}

		hash = sha512.Sum512(hash[:])

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

			hash := sha256.Sum256(entry.value)
			entry.value = hash[:]
			knownKey.sha256[idx] = entry

			return knownKey.siteId, nil
		}

		hash := sha256.Sum256(append(knownKey.key, salt...))
		hash = sha256.Sum256(hash[:])

		if !bytes.Equal(hash[:], key) {
			return 0, errorInvalidKey
		}

		hash = sha256.Sum256(hash[:])

		knownKey.sha256 = append(knownKey.sha256, rollingHash{
			salt:  salt,
			value: hash[:],
		})
	}

	return knownKey.siteId, nil
}
