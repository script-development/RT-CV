package controller

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	. "github.com/stretchr/testify/assert"
)

func TestRouteAuthSeed(t *testing.T) {
	app := newTestingRouter(t)

	_, body := app.MakeRequest(
		routeBuilder.Get,
		"/api/v1/auth/seed",
		TestReqOpts{
			// This route should work without auth
			NoAuth: true,
		},
	)
	Equal(t, fmt.Sprintf(`{"seed":"%s"}`, string(testingServerSeed)), string(body))
}

func TestRouteGetKeyInfo(t *testing.T) {
	app := newTestingRouter(t)

	_, body := app.MakeRequest(
		routeBuilder.Get,
		"/api/v1/auth/keyinfo",
		TestReqOpts{},
	)
	bodyValues := map[string]interface{}{}
	err := json.Unmarshal(body, &bodyValues)
	NoError(t, err)

	for key, value := range bodyValues {
		fmt.Println(key, value)
		NotEqual(t, strings.ToLower(key), "key", "the key should re-appear in the result data")
	}
}
