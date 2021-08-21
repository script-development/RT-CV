package controller

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestSecretRoutes(t *testing.T) {
	app := newTestingRouter()

	contents := `{"key":"value"}`
	valueKey := "test1"
	encryptionKey := "very-secret-key-of-minimal-16-chars"

	// Insert key works
	route := fmt.Sprintf("/v1/scraper/secret/%v/%v", valueKey, encryptionKey)
	res := app.MakeRequest(t, Post, route, TestReqOpts{
		Body: []byte(contents),
	})
	body, err := ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, contents, string(body))

	// Get key works
	res = app.MakeRequest(t, Get, route, TestReqOpts{})
	body, err = ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, contents, string(body))

	// Update the secret
	contents = `{"key":"other value"}`
	res = app.MakeRequest(t, Put, route, TestReqOpts{
		Body: []byte(contents),
	})
	body, err = ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, contents, string(body))

	// Check if we do a get request we recive the updated value
	res = app.MakeRequest(t, Get, route, TestReqOpts{})
	body, err = ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, contents, string(body))

	// Can delete value
	deleteRoute := fmt.Sprintf("/v1/scraper/secret/%v", valueKey)
	res = app.MakeRequest(t, Delete, deleteRoute, TestReqOpts{})
	body, err = ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, `{"status":"ok"}`, string(body))

	// Check if the value is for real deleted
	res = app.MakeRequest(t, Get, route, TestReqOpts{})
	body, err = ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, `{"error":"item not found"}`, string(body))
}
