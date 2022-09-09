package controller

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestStuff(t *testing.T) {
	app := newTestingRouter(t)

	hookDataChan := make(chan string)
	hookWebserver := fiber.New(fiber.Config{DisableStartupMessage: true})
	hookWebserver.Post(`/hook`, func(c *fiber.Ctx) error {
		hookDataChan <- string(c.Body())
		return c.SendString("ok")
	})
	hookWebserverAddr := "127.0.0.1:8989"
	go func() {
		err := hookWebserver.Listen(hookWebserverAddr)
		NoError(t, err)
	}()
	defer hookWebserver.Shutdown()

	// Create a new list profile to match the example cv with
	reqBody, err := json.Marshal(models.Profile{
		Name:         "test list profile",
		Active:       true,
		ListsAllowed: true,
		Zipcodes: []models.ProfileDutchZipcode{{
			From: 9000,
			To:   9999,
		}},
	})
	NoError(t, err)
	res, body := app.MakeRequest(routeBuilder.Post, `/api/v1/profiles`, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	// Create a hook the match can be send to
	reqBody, err = json.Marshal(map[string]string{
		"method": "POST",
		"url":    "http://" + hookWebserverAddr + "/hook",
	})
	NoError(t, err)
	res, body = app.MakeRequest(routeBuilder.Post, `/api/v1/onMatchHooks`, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	// Check if the request was successful
	app.ChangeAuthKey(mock.Key2)
	reqBody, err = json.Marshal(RouteScraperListCVsReq{CVs: []models.CV{*models.ExampleCV()}})
	NoError(t, err)
	res, body = app.MakeRequest(routeBuilder.Post, `/api/v1/scraper/allCvs`, TestReqOpts{Body: reqBody})
	Equal(t, 200, res.StatusCode, string(body))

	// Wait for the hook
	hookData := CVListsHookData{}
	select {
	case data := <-hookDataChan:
		err = json.Unmarshal([]byte(data), &hookData)
		NoError(t, err)
	case <-time.After(time.Second * 2):
		FailNow(t, "exepected hook to be called")
	}

	Len(t, hookData.CVs, 1)
	Len(t, hookData.ProfilesMatchCVs, 1)
	False(t, hookData.IsTest)
}
