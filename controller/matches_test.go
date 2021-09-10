package controller

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestRouteGetMatchesWithinRegion(t *testing.T) {
	from := time.Now().Add(-time.Hour * 24).Format(time.RFC3339)
	to := time.Now().Format(time.RFC3339)

	r := newTestingRouter(t)

	_, body := r.MakeRequest(
		routeBuilder.Get,
		fmt.Sprintf("/api/v1/analytics/matches/period/%s/%s", from, to),
		TestReqOpts{},
	)

	bodyMatches := []models.Match{}
	err := json.Unmarshal(body, &bodyMatches)
	NoError(t, err)

	// There are 2 dummy matches in the database
	Len(t, bodyMatches, 2)
}
