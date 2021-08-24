package controller

import (
	"encoding/json"
	"testing"

	"github.com/script-development/RT-CV/helpers/random"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestApiKeyRoutes(t *testing.T) {
	app := newTestingRouter(t)

	// Get all api keys
	_, res := app.MakeRequest(Get, `/api/v1/keys`, TestReqOpts{})

	// Check if the response contains the api keys inserted in the mock data
	resKeys := []models.APIKey{}
	err := json.Unmarshal(res, &resKeys)
	NoError(t, err)
	Len(t, resKeys, 3) // The mock data contains 3 keys

	// get current keys in the
	allKeysInDB := []models.APIKey{}
	err = app.db.Find(&models.APIKey{}, &allKeysInDB, nil)
	NoError(t, err)

	// Check if the keys in the db matches the repsone
	allKeysInDBJson, err := json.Marshal(allKeysInDB)
	NoError(t, err)
	Equal(t, string(allKeysInDBJson), string(res))

	// Get each key from earlier by id
	for _, listKey := range resKeys {
		keyRoute := `/api/v1/keys/` + listKey.ID.Hex()
		_, res = app.MakeRequest(Get, keyRoute, TestReqOpts{})

		resKey := &models.APIKey{}
		err = json.Unmarshal(res, resKey)
		NoError(t, err)
		Equal(t, listKey.ID.Hex(), resKey.ID.Hex())

		// Delete the key and check if it's really deleted
		// Firstly we count how many document we have before the delete
		keysCountBeforeDeletion := len(resKeys)

		// Send the delete request
		app.MakeRequest(Delete, keyRoute, TestReqOpts{})

		// Count how many keys we have after the deletion
		_, res := app.MakeRequest(Get, `/api/v1/keys`, TestReqOpts{})
		resKeys = []models.APIKey{}
		err = json.Unmarshal(res, &resKeys)
		NoError(t, err)

		Equal(t, keysCountBeforeDeletion-1, len(resKeys))
	}

	// Try to insert key
	randomKey := string(random.GenerateKey())
	trueVal := true
	roles := models.APIKeyRoleScraper
	keyToInsert := apiKeyModifyCreateData{
		Enabled: &trueVal,
		Domains: []string{"example.com"},
		Key:     &randomKey,
		Roles:   &roles,
	}
	body, err := json.Marshal(keyToInsert)
	NoError(t, err)
	_, res = app.MakeRequest(Post, `/api/v1/keys`, TestReqOpts{Body: body})
	resKey := &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	NotNil(t, resKey.ID)
	Equal(t, *keyToInsert.Key, resKey.Key)

	// Check if we can fetch the newly inserted key
	_, res = app.MakeRequest(Get, `/api/v1/keys/`+resKey.ID.Hex(), TestReqOpts{})
	resKey = &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	Equal(t, *keyToInsert.Key, resKey.Key)

	// Try to update the key
	newRandomKey := string(random.GenerateKey())
	_, res = app.MakeRequest(Put, `/api/v1/keys/`+resKey.ID.Hex(), TestReqOpts{Body: []byte(`{"key": "` + newRandomKey + `"}`)})
	resKey = &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	Equal(t, newRandomKey, resKey.Key)

	// check if the key was updated
	_, res = app.MakeRequest(Get, `/api/v1/keys/`+resKey.ID.Hex(), TestReqOpts{})
	resKey = &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	Equal(t, newRandomKey, resKey.Key)
}
