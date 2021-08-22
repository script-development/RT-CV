package controller

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestRouteControlReloadProfiles(t *testing.T) {
	app := newTestingRouter()

	_, body := app.MakeRequest(t, Get, "/v1/control/reloadProfiles", TestReqOpts{})
	Equal(t, `{"status":"ok"}`, string(body))
}
