package controller

import (
	"fmt"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	. "github.com/stretchr/testify/assert"
)

func TestSecretRoutes(t *testing.T) {
	app := newTestingRouter(t)

	contents := `{"key":"value"}`
	valueKey := "test1"
	encryptionKey := "very-secret-key-of-minimal-16-chars"

	// Insert key works
	route := fmt.Sprintf("/api/v1/secrets/myKey/%v/%v", valueKey, encryptionKey)
	_, body := app.MakeRequest(routeBuilder.Post, route, TestReqOpts{
		Body: []byte(fmt.Sprintf(`{"value": %s}`, contents)),
	})
	Equal(t, contents, string(body))

	// Get key works
	_, body = app.MakeRequest(routeBuilder.Get, route, TestReqOpts{})
	Equal(t, contents, string(body))

	// Update the secret
	contents = `{"key":"other value"}`
	_, body = app.MakeRequest(routeBuilder.Put, route, TestReqOpts{
		Body: []byte(fmt.Sprintf(`{"value": %s}`, contents)),
	})
	Equal(t, contents, string(body))

	// Check if we do a get request we recive the updated value
	_, body = app.MakeRequest(routeBuilder.Get, route, TestReqOpts{})
	Equal(t, contents, string(body))

	// Can delete value
	deleteRoute := fmt.Sprintf("/api/v1/secrets/myKey/%v", valueKey)
	_, body = app.MakeRequest(routeBuilder.Delete, deleteRoute, TestReqOpts{})
	Equal(t, `{"status":"ok"}`, string(body))

	// Check if the value is for real deleted
	_, body = app.MakeRequest(routeBuilder.Get, route, TestReqOpts{})
	Equal(t, `{"error":"item not found"}`, string(body))
}
