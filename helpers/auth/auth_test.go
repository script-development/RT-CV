package auth

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accessor struct {
	rollingKey      []byte
	keyBytes        []byte
	saltBytes       []byte
	keyAndSaltBytes []byte // = keyBytes + saltBytes
	keyId           string
	baseSeed        []byte
}

func newAccessorHelper(keyId primitive.ObjectID, key, salt string) *accessor {
	keyBytes := []byte(key)
	saltBytes := []byte(salt)
	keyandSaltBytes := append(keyBytes, saltBytes...)

	h := sha512.Sum512(keyandSaltBytes)
	return &accessor{
		rollingKey:      h[:],
		keyBytes:        keyBytes,
		saltBytes:       saltBytes,
		keyAndSaltBytes: keyandSaltBytes,
		keyId:           keyId.Hex(),
		baseSeed:        []byte("unsafe-testing-seed"),
	}
}

func (a *accessor) key() []byte {
	newRollingKey := sha512.Sum512(append(append(a.rollingKey, a.keyBytes...), a.saltBytes...))
	a.rollingKey = newRollingKey[:]

	src := bytes.Join([][]byte{
		[]byte("sha512:" + a.keyId),
		a.saltBytes,
		[]byte(hex.EncodeToString(a.rollingKey)),
	}, []byte(":"))

	return []byte("Basic " + base64.RawStdEncoding.EncodeToString(src))
}

func TestAuthenticate(t *testing.T) {
	key1ID := primitive.NewObjectID()
	key2ID := primitive.NewObjectID()
	key3ID := primitive.NewObjectID()

	auth := New([]models.ApiKey{
		{
			M:       dbInterfaces.M{ID: key1ID},
			Enabled: true,
			Domains: []string{"a", "b"},
			Key:     "abc",
			Roles:   models.ApiKeyRoleScraper,
		},
		{
			M:       dbInterfaces.M{ID: key2ID},
			Enabled: true,
			Domains: []string{"c", "d"},
			Key:     "def",
			Roles:   models.ApiKeyRoleInformationObtainer,
		},
		{
			M:       dbInterfaces.M{ID: key3ID},
			Enabled: true,
			Domains: []string{"e", "f"},
			Key:     "ghi",
			Roles:   models.ApiKeyRoleController,
		},
	}, []byte("unsafe-testing-seed"))

	// No key provided
	_, _, err := auth.Authenticate([]byte{})
	assert.Error(t, err)

	// First time key usage
	site1KeySaltFoo := newAccessorHelper(key1ID, "abc", "foo")
	key := site1KeySaltFoo.key()
	site, _, err := auth.Authenticate(key)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, site.Domains)

	// Using the same key twice should yield an error
	_, _, err = auth.Authenticate(key)
	assert.Error(t, err)

	// Generating a new key should work
	site, _, err = auth.Authenticate(site1KeySaltFoo.key())
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, site.Domains)

	// Creating a new key for the same site with a diffrent salt should work
	site1KeySaltBar := newAccessorHelper(key1ID, "abc", "bar")
	site, _, err = auth.Authenticate(site1KeySaltBar.key())
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, site.Domains)

	// Generating a new key using the second key should work
	_, _, err = auth.Authenticate(site1KeySaltBar.key())
	assert.NoError(t, err)

	// The first key created should also still work
	_, _, err = auth.Authenticate(site1KeySaltFoo.key())
	assert.NoError(t, err)

	// Using the wrong input key should fail
	_, _, err = auth.Authenticate(newAccessorHelper(key1ID, "this-is-the-wrong-key", "baz").key())
	assert.Error(t, err)

	// Authenticating a diffrent site should work
	site2KeySaltFoo := newAccessorHelper(key2ID, "def", "foo")
	site, _, err = auth.Authenticate(site2KeySaltFoo.key())
	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, site.Domains)

	// Using another key id's key should fail
	site1WithKeyFrom2 := newAccessorHelper(key1ID, "def", "foobar")
	_, _, err = auth.Authenticate(site1WithKeyFrom2.key())
	assert.Error(t, err)
}
