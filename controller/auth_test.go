package controller

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestRouteAuthSeed(t *testing.T) {
	app := newTestingRouter()

	_, body := app.MakeRequest(
		t,
		Get,
		"/v1/auth/seed",
		TestReqOpts{
			// This route should work without auth
			NoAuth: true,
		},
	)
	Equal(t, fmt.Sprintf(`{"seed":"%s"}`, string(testingServerSeed)), string(body))
}
