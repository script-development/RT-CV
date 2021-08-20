package auth

import (
	"testing"

	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthenticate(t *testing.T) {
	key1ID := primitive.NewObjectID()
	key2ID := primitive.NewObjectID()
	key3ID := primitive.NewObjectID()

	authSeed := []byte("unsafe-testing-seed")

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
	}, authSeed)

	// No key provided
	_, _, err := auth.Authenticate([]byte{})
	assert.Error(t, err)

	// First time key usage
	site1KeySaltFoo := NewAccessorHelper(key1ID, "abc", "foo", authSeed)
	key := site1KeySaltFoo.Key()
	site, _, err := auth.Authenticate(key)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, site.Domains)

	// Using the same key twice should yield an error
	_, _, err = auth.Authenticate(key)
	assert.Error(t, err)

	// Generating a new key should work
	site, _, err = auth.Authenticate(site1KeySaltFoo.Key())
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, site.Domains)

	// Creating a new key for the same site with a diffrent salt should work
	site1KeySaltBar := NewAccessorHelper(key1ID, "abc", "bar", authSeed)
	site, _, err = auth.Authenticate(site1KeySaltBar.Key())
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, site.Domains)

	// Generating a new key using the second key should work
	_, _, err = auth.Authenticate(site1KeySaltBar.Key())
	assert.NoError(t, err)

	// The first key created should also still work
	_, _, err = auth.Authenticate(site1KeySaltFoo.Key())
	assert.NoError(t, err)

	// Using the wrong input key should fail
	site1WithWrongKey := NewAccessorHelper(key1ID, "this-is-a-wrong-key", "baz", authSeed)
	_, _, err = auth.Authenticate(site1WithWrongKey.Key())
	assert.Error(t, err)

	// Authenticating a diffrent site should work
	site2KeySaltFoo := NewAccessorHelper(key2ID, "def", "foo", authSeed)
	site, _, err = auth.Authenticate(site2KeySaltFoo.Key())
	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, site.Domains)

	// Using another key id's key should fail
	site1WithKeyFrom2 := NewAccessorHelper(key1ID, "def", "foobar", authSeed)
	_, _, err = auth.Authenticate(site1WithKeyFrom2.Key())
	assert.Error(t, err)
}
