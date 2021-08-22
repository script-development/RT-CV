package controller

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestRouteAuthSeed(t *testing.T) {
	app := newTestingRouter(t)

	_, body := app.MakeRequest(
		Get,
		"/api/v1/auth/seed",
		TestReqOpts{
			// This route should work without auth
			NoAuth: true,
		},
	)
	Equal(t, fmt.Sprintf(`{"seed":"%s"}`, string(testingServerSeed)), string(body))
}
