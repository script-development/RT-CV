package controller

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	. "github.com/stretchr/testify/assert"
)

func TestRouteGetKeyInfo(t *testing.T) {
	app := newTestingRouter(t)

	_, body := app.MakeRequest(
		routeBuilder.Get,
		"/api/v1/auth/keyinfo",
		TestReqOpts{},
	)
	bodyValues := map[string]any{}
	err := json.Unmarshal(body, &bodyValues)
	NoError(t, err)

	for key := range bodyValues {
		NotEqual(t, strings.ToLower(key), "key", "the key property should re-appear in the result data")
	}
}
