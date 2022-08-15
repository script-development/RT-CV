package controller

import (
	"encoding/json"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestSecretRoutes(t *testing.T) {
	app := newTestingRouter(t)

	contents := `{"key":"value"}`
	valueKey := "test1"
	encryptionKey := "very-secret-key-of-minimal-16-chars"

	createBody := RouteUpdateOrCreateSecret{
		Value:          json.RawMessage(contents),
		ValueStructure: models.SecretValueStructureFree,
		Description:    "",
		EncryptionKey:  encryptionKey,
	}
	createBodyJSON, _ := json.Marshal(createBody)

	// Insert key works
	route := "/api/v1/secrets/myKey/" + valueKey
	_, body := app.MakeRequest(routeBuilder.Put, route, TestReqOpts{
		Body: createBodyJSON,
	})
	Equal(t, contents, string(body))

	// Get key works
	_, body = app.MakeRequest(routeBuilder.Get, route+"/"+encryptionKey, TestReqOpts{})
	Equal(t, contents, string(body))

	// Update the secret
	contents = `{"key":"other value"}`
	createBody.Value = json.RawMessage(contents)
	createBodyJSON, _ = json.Marshal(createBody)

	_, body = app.MakeRequest(routeBuilder.Put, route, TestReqOpts{
		Body: []byte(createBodyJSON),
	})
	Equal(t, contents, string(body))

	// Check if we do a get request we recive the updated value
	_, body = app.MakeRequest(routeBuilder.Get, route+"/"+encryptionKey, TestReqOpts{})
	Equal(t, contents, string(body))

	// Can delete value
	_, body = app.MakeRequest(routeBuilder.Delete, route, TestReqOpts{})
	Equal(t, `{"status":"ok"}`, string(body))

	// Check if the value is for real deleted
	_, body = app.MakeRequest(routeBuilder.Get, route+"/"+encryptionKey, TestReqOpts{})
	Equal(t, `{"error":"item not found"}`, string(body))
}
