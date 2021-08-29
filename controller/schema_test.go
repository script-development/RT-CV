package controller

import (
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
)

func TestRouteGetCvSchema(t *testing.T) {
	newTestingRouter(t).MakeRequest(
		routeBuilder.Get,
		"/api/v1/schema/cv",
		TestReqOpts{NoAuth: true},
	)
}
