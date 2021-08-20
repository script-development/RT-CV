package controller

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestRouteAuthSeed(t *testing.T) {
	app := newTestingRouter()

	res := app.MakeRequest(
		t,
		Get,
		"/v1/auth/seed",
		TestReqOpts{
			// This route should work without auth
			NoAuth: true,
		},
	)

	body, err := ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, fmt.Sprintf(`{"seed":"%s"}`, string(testingServerSeed)), string(body))
}
