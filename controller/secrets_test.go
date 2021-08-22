package controller

import (
	"fmt"
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
	_, body := app.MakeRequest(t, Post, route, TestReqOpts{
		Body: []byte(contents),
	})
	Equal(t, contents, string(body))

	// Get key works
	_, body = app.MakeRequest(t, Get, route, TestReqOpts{})
	Equal(t, contents, string(body))

	// Update the secret
	contents = `{"key":"other value"}`
	_, body = app.MakeRequest(t, Put, route, TestReqOpts{
		Body: []byte(contents),
	})
	Equal(t, contents, string(body))

	// Check if we do a get request we recive the updated value
	_, body = app.MakeRequest(t, Get, route, TestReqOpts{})
	Equal(t, contents, string(body))

	// Can delete value
	deleteRoute := fmt.Sprintf("/v1/scraper/secret/%v", valueKey)
	_, body = app.MakeRequest(t, Delete, deleteRoute, TestReqOpts{})
	Equal(t, `{"status":"ok"}`, string(body))

	// Check if the value is for real deleted
	_, body = app.MakeRequest(t, Get, route, TestReqOpts{})
	Equal(t, `{"error":"item not found"}`, string(body))
}
