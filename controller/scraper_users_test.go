package controller

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
	"golang.org/x/crypto/nacl/box"
)

func TestScraperUsers(t *testing.T) {
	r := newTestingRouter(t)
	r.ChangeAuthKey(mock.Key1)

	path := "/api/v1/scraperUsers/" + mock.Key1.ID.Hex()
	res, body := r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	scraperUsers := models.ScraperLoginUsers{}
	err := json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 0, len(scraperUsers.Users))

	// ---
	// Add public key used for encryption

	pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)
	NoError(t, err)

	testDecrypt := func(cipherBase64, expected string) {
		cipherBytes, err := base64.StdEncoding.DecodeString(cipherBase64)
		NoError(t, err)
		decrypted, ok := box.OpenAnonymous(nil, cipherBytes, pubKey, privKey)
		True(t, ok)
		Equal(t, expected, string(decrypted[32:]))
	}

	pubKeyStr := base64.StdEncoding.EncodeToString(pubKey[:])
	res, body = r.MakeRequest(routeBuilder.Patch, path+"/setPublicKey", TestReqOpts{Body: []byte(`{"publicKey":"` + pubKeyStr + `"}`)})
	Equal(t, 200, res.StatusCode, string(body))

	scraperUsers = models.ScraperLoginUsers{}
	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, pubKeyStr, scraperUsers.ScraperPubKey)

	// ---
	// Add user

	reqBody := []byte(`{"username": "username", "password": "password"}`)
	res, body = r.MakeRequest(routeBuilder.Patch, path, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 1, len(scraperUsers.Users))

	// ---
	// Check if adding a user was successful

	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 1, len(scraperUsers.Users))
	usr := scraperUsers.Users[0]
	Equal(t, "username", usr.Username)
	testDecrypt(usr.EncryptedPassword, "password")

	// ---
	// Check if requesting the user using a non scraper key hides the passwords

	r.ChangeAuthKey(mock.Key3)

	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	scraperUsers = models.ScraperLoginUsers{}
	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))

	Equal(t, 1, len(scraperUsers.Users))
	Equal(t, models.ScraperLoginUser{Username: "username"}, scraperUsers.Users[0], string(body))

	r.ChangeAuthKey(mock.Key1)

	// ---
	// Add another user

	reqBody = []byte(`{"username": "username2", "password": "password2"}`)
	res, body = r.MakeRequest(routeBuilder.Patch, path, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 2, len(scraperUsers.Users))

	// ---
	// Check if adding a user was successful

	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 2, len(scraperUsers.Users))
	Equal(t, "username", scraperUsers.Users[0].Username)
	Equal(t, "username2", scraperUsers.Users[1].Username)

	// ---
	// Update a user

	reqBody = []byte(`{"username": "username", "password": "updated password"}`)
	res, body = r.MakeRequest(routeBuilder.Patch, path, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 2, len(scraperUsers.Users))

	// ---
	// Delete a user

	reqBody = []byte(`{"username": "username2"}`)
	res, body = r.MakeRequest(routeBuilder.Delete, path, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 1, len(scraperUsers.Users), string(body))

	// ---
	// Check if the update and delete were successful

	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 1, len(scraperUsers.Users))
	usr = scraperUsers.Users[0]
	Equal(t, "username", usr.Username)
	testDecrypt(usr.EncryptedPassword, "updated password")

	// ---
	// Check if requesting the login users of another scraper returns an empty array

	path = "/api/v1/scraperUsers/" + mock.Key2.ID.Hex()
	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 0, len(scraperUsers.Users), string(body))
}
