package controller

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestScraperUsers(t *testing.T) {
	r := newTestingRouter(t)

	path := "/api/v1/scraperUsers/" + mock.Key1.ID.Hex()
	res, body := r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	scraperUsers := models.ScraperLoginUsers{}
	err := json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 0, len(scraperUsers.Users))

	// ---
	// Add user

	reqBody := []byte(`{"username": "username", "password": "password"}`)
	res, body = r.MakeRequest(routeBuilder.Patch, path, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 1, len(scraperUsers.Users))

	// ---
	// Check if added user was successful

	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 1, len(scraperUsers.Users))

	// ---
	// Add another user

	reqBody = []byte(`{"username": "username2", "password": "password2"}`)
	res, body = r.MakeRequest(routeBuilder.Patch, path, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 2, len(scraperUsers.Users))

	// ---
	// Check if added user was successful

	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 2, len(scraperUsers.Users))
	Equal(t, scraperUsers.Users, []models.ScraperLoginUser{
		{Username: "username", Password: "password"},
		{Username: "username2", Password: "password2"},
	})

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
	Equal(t, scraperUsers.Users[0], models.ScraperLoginUser{
		Username: "username",
		Password: "updated password",
	})

	// ---
	// Check if requesting the login users of another scraper returns an empty array

	path = "/api/v1/scraperUsers/" + mock.Key2.ID.Hex()
	res, body = r.MakeRequest(routeBuilder.Get, path, TestReqOpts{})
	Equal(t, 200, res.StatusCode, string(body))

	fmt.Println(mock.Key2.ID.Hex())

	err = json.Unmarshal(body, &scraperUsers)
	NoError(t, err, string(body))
	Equal(t, 0, len(scraperUsers.Users), string(body))
}
