package controller

import (
	"encoding/json"
	"testing"

	"github.com/script-development/RT-CV/helpers/random"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestApiKeyRoutes(t *testing.T) {
	app := newTestingRouter(t)

	// Get all api keys
	_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/keys`, TestReqOpts{})

	// Check if the response contains the api keys inserted in the mock data
	resKeys := []models.APIKey{}
	err := json.Unmarshal(res, &resKeys)
	NoError(t, err)
	Len(t, resKeys, 4) // The mock data contains 4 keys

	// get current keys in the
	allKeysInDB := []models.APIKey{}
	err = app.db.Find(&models.APIKey{}, &allKeysInDB, nil)
	NoError(t, err)

	// Check if the keys in the db matches the repsone
	allKeysInDBJson, err := json.Marshal(allKeysInDB)
	NoError(t, err)
	Equal(t, string(allKeysInDBJson), string(res))

	// Get scraper keys
	_, scraperKeysResp := app.MakeRequest(routeBuilder.Get, `/api/v1/keys/scrapers`, TestReqOpts{})
	scraperKeys := []models.APIKey{}
	err = json.Unmarshal(scraperKeysResp, &scraperKeys)
	NoError(t, err)
	Len(t, scraperKeys, 1) // The mock data contains one scraper keys

	// Get each key from earlier by id
	for _, listKey := range resKeys {
		if listKey.ID.Hex() == mock.Key1.ID.Hex() {
			continue
		}

		keyRoute := `/api/v1/keys/` + listKey.ID.Hex()
		_, res = app.MakeRequest(routeBuilder.Get, keyRoute, TestReqOpts{})

		resKey := &models.APIKey{}
		err = json.Unmarshal(res, resKey)
		NoError(t, err)
		Equal(t, listKey.ID.Hex(), resKey.ID.Hex())

		// Delete the key and check if it's really deleted
		// Firstly we count how many document we have before the delete
		keysCountBeforeDeletion := len(resKeys)

		// Send the delete request
		app.MakeRequest(routeBuilder.Delete, keyRoute, TestReqOpts{})

		// Count how many keys we have after the deletion
		_, res := app.MakeRequest(routeBuilder.Get, `/api/v1/keys`, TestReqOpts{})
		resKeys = []models.APIKey{}
		err = json.Unmarshal(res, &resKeys)
		NoError(t, err)

		if listKey.System {
			// System keys cannot be removed
			Equal(t, keysCountBeforeDeletion, len(resKeys))
		} else {
			Equal(t, keysCountBeforeDeletion-1, len(resKeys))
		}
	}

	// Try to insert key
	randomKey := string(random.GenerateKey())
	name := "test example key"
	trueVal := true
	roles := models.APIKeyRoleScraper
	keyToInsert := apiKeyModifyCreateData{
		Enabled: &trueVal,
		Name:    &name,
		Domains: []string{"example.com"},
		Key:     &randomKey,
		Roles:   &roles,
	}
	body, err := json.Marshal(keyToInsert)
	NoError(t, err)
	_, res = app.MakeRequest(routeBuilder.Post, `/api/v1/keys`, TestReqOpts{Body: body})
	resKey := &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	NotNil(t, resKey.ID)
	Equal(t, *keyToInsert.Key, resKey.Key)

	// Check if we can fetch the newly inserted key
	_, res = app.MakeRequest(routeBuilder.Get, `/api/v1/keys/`+resKey.ID.Hex(), TestReqOpts{})
	resKey = &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	Equal(t, *keyToInsert.Key, resKey.Key)

	// Try to update the key
	newRandomKey := string(random.GenerateKey())
	_, res = app.MakeRequest(
		routeBuilder.Put,
		`/api/v1/keys/`+resKey.ID.Hex(),
		TestReqOpts{Body: []byte(`{"key": "` + newRandomKey + `"}`)},
	)
	resKey = &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	Equal(t, newRandomKey, resKey.Key)

	// check if the key was updated
	_, res = app.MakeRequest(routeBuilder.Get, `/api/v1/keys/`+resKey.ID.Hex(), TestReqOpts{})
	resKey = &models.APIKey{}
	err = json.Unmarshal(res, resKey)
	NoError(t, err)
	Equal(t, newRandomKey, resKey.Key)
}
