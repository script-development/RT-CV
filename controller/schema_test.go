package controller

import (
	"testing"
)

func TestRouteGetCvSchema(t *testing.T) {
	newTestingRouter(t).MakeRequest(
		Get,
		"/api/v1/schema/cv",
		TestReqOpts{NoAuth: true},
	)
}
