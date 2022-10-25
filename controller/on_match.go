package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeGetOnMatchHooks = routeBuilder.R{
	Description: "Get the entries for what to do on a match",
	Res:         []models.OnMatchHook{},
	Fn: func(c *fiber.Ctx) error {
		results := []models.OnMatchHook{}
		err := ctx.Get(c).DBConn.Find(&models.OnMatchHook{}, &results, nil)
		if err != nil {
			return err
		}
		return c.JSON(results)
	},
}

// CreateOrUpdateOnMatchHookRequestData contains the post data for creating and modifiying a OnMatchHook
type CreateOrUpdateOnMatchHookRequestData struct {
	Disabled   *bool           `json:"disabled"`
	Method     *string         `json:"method"`
	URL        *string         `json:"url"`
	AddHeaders []models.Header `json:"addHeaders"`
}

func (data *CreateOrUpdateOnMatchHookRequestData) applyToHook(hook *models.OnMatchHook, isCreate bool) error {
	if data.Disabled != nil {
		hook.Disabled = *data.Disabled
	}

	if data.Method != nil {
		method := strings.ToUpper(*data.Method)
		allowedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
		methodAllowed := false
		for _, allowedMethod := range allowedMethods {
			if method == allowedMethod {
				methodAllowed = true
				break
			}
		}
		if !methodAllowed {
			return errors.New("method should be one of GET, POST, PUT, PATCH, DELETE")
		}
		hook.Method = method
	} else if isCreate {
		hook.Method = "GET"
	}

	if data.URL != nil {
		url := *data.URL
		if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
			return errors.New("url should start with http(s)://")
		}
		hook.URL = url
	} else if isCreate {
		return errors.New("url is required")
	}

	if data.AddHeaders != nil {
		hook.AddHeaders = []models.Header{}
		for _, extraHeader := range data.AddHeaders {
			if extraHeader.Key != "" {
				hook.AddHeaders = append(hook.AddHeaders, extraHeader)
			}
		}
	} else if isCreate {
		hook.AddHeaders = []models.Header{}
	}

	return nil
}

var routeCreateOnMatchHooks = routeBuilder.R{
	Description: "Set a on match hook, called when a cv is matched with one or more profiles or when a list of cvs is scanned and matched to list profiles",
	Body:        CreateOrUpdateOnMatchHookRequestData{},
	Res:         models.OnMatchHook{},
	Fn: func(c *fiber.Ctx) error {
		body := CreateOrUpdateOnMatchHookRequestData{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}
		ctx := ctx.Get(c)

		hook := models.OnMatchHook{
			M:     db.NewM(),
			KeyID: ctx.Key.ID,
		}
		err = body.applyToHook(&hook, true)
		if err != nil {
			return err
		}

		err = ctx.DBConn.Insert(&hook)
		if err != nil {
			return err
		}

		return c.JSON(body)
	},
}

var routeDeleteOnMatchHook = routeBuilder.R{
	Description: "Delete a on match hook",
	Res:         models.OnMatchHook{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)
		err := ctx.DBConn.DeleteByID(&models.OnMatchHook{}, ctx.OnMatchHook.ID)
		if err != nil {
			return err
		}
		return c.JSON(ctx.OnMatchHook)
	},
}

var routeUpdateOnMatchHook = routeBuilder.R{
	Description: "Update a on match hook",
	Body:        CreateOrUpdateOnMatchHookRequestData{},
	Res:         models.OnMatchHook{},
	Fn: func(c *fiber.Ctx) error {
		body := CreateOrUpdateOnMatchHookRequestData{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}
		ctx := ctx.Get(c)

		err = body.applyToHook(ctx.OnMatchHook, false)
		if err != nil {
			return err
		}

		err = ctx.DBConn.UpdateByID(ctx.OnMatchHook)
		if err != nil {
			return err
		}
		return c.JSON(ctx.OnMatchHook)
	},
}

// RouteTestOnMatchHookExplainSendToHook explains what is send to the hook
type RouteTestOnMatchHookExplainSendToHook struct {
	DataSendToHook    HookMatchedCVData `json:"dataSendToHook"`
	HeadersSendToHook http.Header       `json:"headers"`
}

var routeTestOnMatchHook = routeBuilder.R{
	Description: "Test a on match hook",
	Res:         RouteTestOnMatchHookExplainSendToHook{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)
		cv := *models.ExampleCV()

		yearsSinceWork := 3

		dummyData := HookMatchedCVData{
			MatchedProfiles: []match.FoundMatch{{
				Matches: models.Match{
					RequestID:         ctx.RequestID,
					ProfileID:         mock.Profile1.ID,
					KeyID:             mock.Key1.ID,
					When:              jsonHelpers.RFC3339Nano(time.Now()),
					ReferenceNr:       cv.ReferenceNumber,
					Debug:             false,
					YearsSinceWork:    &yearsSinceWork,
					DesiredProfession: &mock.Profile1.DesiredProfessions[0].Name,
					ZipCode:           &models.ProfileDutchZipcode{},
				},
				Profile: *mock.Profile1,
			}},
			CV:      cv,
			KeyID:   mock.Key1.ID,
			KeyName: "example.com",
			IsTest:  true,
		}

		dummyDataAsJSON, err := json.Marshal(dummyData)
		if err != nil {
			return err
		}

		headers, err := ctx.OnMatchHook.CallWithRetry(bytes.NewReader(dummyDataAsJSON), models.DataKindMatch)
		if err != nil {
			return err
		}

		return c.JSON(RouteTestOnMatchHookExplainSendToHook{DataSendToHook: dummyData, HeadersSendToHook: headers})
	},
}

// RouteTestOnListHookExplainSendToHook explains what is send to the hook
type RouteTestOnListHookExplainSendToHook struct {
	DataSendToHook    CVListsHookData `json:"dataSendToHook"`
	HeadersSendToHook http.Header     `json:"headers"`
}

var routeTestOnListHook = routeBuilder.R{
	Description: "Test a on list hook",
	Res:         RouteTestOnListHookExplainSendToHook{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)
		cv := *models.ExampleCV()

		dummyData := CVListsHookData{
			CVs: map[string]models.CV{
				cv.ReferenceNumber: cv,
			},
			ProfilesMatchCVs: map[primitive.ObjectID][]string{
				mock.Profile1.ID: {cv.ReferenceNumber},
			},
			KeyID:   mock.Key1.ID,
			KeyName: "example.com",
			IsTest:  true,
		}

		dummyDataAsJSON, err := json.Marshal(dummyData)
		if err != nil {
			return err
		}

		headers, err := ctx.OnMatchHook.CallWithRetry(bytes.NewReader(dummyDataAsJSON), models.DataKindList)
		if err != nil {
			return err
		}

		return c.JSON(RouteTestOnListHookExplainSendToHook{DataSendToHook: dummyData, HeadersSendToHook: headers})
	},
}

func middlewareBindHook() routeBuilder.M {
	return routeBuilder.M{
		Fn: func(c *fiber.Ctx) error {
			hookParam := c.Params(`hookID`)
			hookID, err := primitive.ObjectIDFromHex(hookParam)
			if err != nil {
				return err
			}
			ctx := ctx.Get(c)
			hook := models.OnMatchHook{}
			query := bson.M{"_id": hookID}
			args := db.FindOptions{NoDefaultFilters: true}
			err = ctx.DBConn.FindOne(&hook, query, args)
			if err != nil {
				return err
			}

			ctx.OnMatchHook = &hook

			return c.Next()
		},
	}
}
