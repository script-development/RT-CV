package auth

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/script-development/RT-CV/models"
	"github.com/stretchr/testify/assert"
)

type accessor struct {
	rollingKey []byte
	keyBytes   []byte
	saltBytes  []byte
	keyId      int
}

func newAccessorHelper(keyId int, key, salt string) *accessor {
	h := sha512.Sum512([]byte(key + salt))
	return &accessor{
		rollingKey: h[:],
		keyBytes:   []byte(key),
		saltBytes:  []byte(salt),
		keyId:      keyId,
	}
}

func (a *accessor) key() []byte {
	newRollingKey := sha512.Sum512(append(append(a.rollingKey, a.keyBytes...), a.saltBytes...))
	a.rollingKey = newRollingKey[:]

	src := bytes.Join([][]byte{
		[]byte("sha512:" + strconv.Itoa(a.keyId)),
		a.saltBytes,
		[]byte(hex.EncodeToString(a.rollingKey)),
	}, []byte(":"))

	return []byte("Basic " + base64.RawStdEncoding.EncodeToString(src))
}

func TestAuthenticate(t *testing.T) {
	auth := New([]models.ApiKey{
		{
			ID:      1,
			Enabled: true,
			SiteId:  4,
			Key:     "abc",
			Roles:   models.ApiKeyRoleScraper,
		},
		{
			ID:      2,
			Enabled: true,
			SiteId:  5,
			Key:     "def",
			Roles:   models.ApiKeyRoleInformationObtainer,
		},
		{
			ID:      3,
			Enabled: true,
			SiteId:  6,
			Key:     "ghi",
			Roles:   models.ApiKeyRoleController,
		},
	})

	// No key provided
	_, err := auth.Authenticate([]byte{})
	assert.Error(t, err)

	// First time key usage
	site1KeySaltFoo := newAccessorHelper(1, "abc", "foo")
	key := site1KeySaltFoo.key()
	siteId, err := auth.Authenticate(key)
	assert.NoError(t, err)
	assert.Equal(t, 4, siteId)

	// Using the same key twice should yield an error
	_, err = auth.Authenticate(key)
	assert.Error(t, err)

	// Generating a new key should work
	siteId, err = auth.Authenticate(site1KeySaltFoo.key())
	assert.NoError(t, err)
	assert.Equal(t, 4, siteId)

	// Creating a new key for the same site with a diffrent salt should work
	site1KeySaltBar := newAccessorHelper(1, "abc", "bar")
	siteId, err = auth.Authenticate(site1KeySaltBar.key())
	assert.NoError(t, err)
	assert.Equal(t, 4, siteId)

	// Generating a new key using the second key should work
	_, err = auth.Authenticate(site1KeySaltBar.key())
	assert.NoError(t, err)

	// The first key created should also still work
	_, err = auth.Authenticate(site1KeySaltFoo.key())
	assert.NoError(t, err)

	// Using the wrong input key should fail
	_, err = auth.Authenticate(newAccessorHelper(1, "this-is-the-wrong-key", "baz").key())
	assert.Error(t, err)

	// Authenticating a diffrent site should work
	site2KeySaltFoo := newAccessorHelper(2, "def", "foo")
	siteId, err = auth.Authenticate(site2KeySaltFoo.key())
	assert.NoError(t, err)
	assert.Equal(t, 5, siteId)

	// Using another key id's key should fail
	site1WithKeyFrom2 := newAccessorHelper(1, "def", "foobar")
	_, err = auth.Authenticate(site1WithKeyFrom2.key())
	assert.Error(t, err)
}
