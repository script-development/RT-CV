package controller

import (
	"encoding/json"
	"testing"

	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestOnMatchHooks(t *testing.T) {
	r := newTestingRouter(t)

	// Create a new hook
	createHookBody, err := json.Marshal(models.OnMatchHook{
		URL:    "https://localhost",
		Method: "get",
		AddHeaders: []models.Header{{
			Key:   "X-Test",
			Value: []string{"test"},
		}},
	})
	NoError(t, err)
	r.MakeRequest(routeBuilder.Post, "/api/v1/onMatchHooks", TestReqOpts{Body: createHookBody})

	// Check if the just inserted hook is actually inserted
	res, body := r.MakeRequest(routeBuilder.Get, "/api/v1/onMatchHooks", TestReqOpts{})
	Equal(t, res.StatusCode, 200, string(body))
	hooks := []models.OnMatchHook{}
	err = json.Unmarshal(body, &hooks)
	NoError(t, err)
	Len(t, hooks, 1)

	firstHook := hooks[0]
	Equal(t, firstHook.URL, "https://localhost")
	Equal(t, firstHook.Method, "GET")
	Len(t, firstHook.AddHeaders, 1)

	firstheader := firstHook.AddHeaders[0]
	Equal(t, firstheader.Key, "X-Test")
	Equal(t, firstheader.Value[0], "test")

	// Delete the just inserted hook
	res, body = r.MakeRequest(routeBuilder.Delete, "/api/v1/onMatchHooks/"+firstHook.ID.Hex(), TestReqOpts{})
	Equal(t, res.StatusCode, 200, string(body))

	// Check if it's deleted
	res, body = r.MakeRequest(routeBuilder.Get, "/api/v1/onMatchHooks", TestReqOpts{})
	Equal(t, res.StatusCode, 200, string(body))
	hooks = []models.OnMatchHook{}
	err = json.Unmarshal(body, &hooks)
	NoError(t, err)
	Len(t, hooks, 0)
}
