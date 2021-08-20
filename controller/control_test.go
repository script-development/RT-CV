package controller

import (
	"io/ioutil"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestRouteControlReloadProfiles(t *testing.T) {
	app := newTestingRouter()

	res := app.MakeRequest(t, Get, "/v1/control/reloadProfiles", TestReqOpts{})

	body, err := ioutil.ReadAll(res.Body)
	NoError(t, err)
	bodyString := string(body)
	Equal(t, `{"status":"ok"}`, bodyString)
}
