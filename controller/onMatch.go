package controller

import (
	"bytes"
	"encoding/json"
	"errors"
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
		db := ctx.GetDbConn(c)
		results := []models.OnMatchHook{}
		err := db.Find(&models.OnMatchHook{}, &results, nil)
		if err != nil {
			return err
		}
		return c.JSON(results)
	},
}

var routeCreateOnMatchHooks = routeBuilder.R{
	Description: "Set a on match hook, called when a cv is matched with one or more profiles\n\n**The id and keyId fields can be left empty**",
	Body:        models.OnMatchHook{},
	Res:         models.OnMatchHook{},
	Fn: func(c *fiber.Ctx) error {
		body := models.OnMatchHook{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		body.ID = primitive.NewObjectID()
		if !strings.HasPrefix(body.URL, "https://") && !strings.HasPrefix(body.URL, "http://") {
			return errors.New("url should start with http(s)://")
		}

		body.KeyID = ctx.GetKey(c).ID

		body.Method = strings.ToUpper(body.Method)
		allowedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
		methodAllowed := false
		for _, allowedMethod := range allowedMethods {
			if body.Method == allowedMethod {
				methodAllowed = true
				break
			}
		}
		if !methodAllowed {
			return errors.New("method should be one of GET, POST, PUT, PATCH, DELETE")
		}

		err = ctx.GetDbConn(c).Insert(&body)
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
		onMatchHook := ctx.GetOnMatchHook(c)
		err := ctx.GetDbConn(c).DeleteByID(onMatchHook)
		if err != nil {
			return err
		}
		return c.JSON(onMatchHook)
	},
}

// ExplainDataSendToHook tells what data is send to the hook
type ExplainDataSendToHook struct {
	DataSendToHook DataSendToHook `json:"dataSendToHook"`
}

var routeTestOnMatchHook = routeBuilder.R{
	Description: "Test a on match hook",
	Res:         ExplainDataSendToHook{},
	Fn: func(c *fiber.Ctx) error {
		requestID := ctx.GetRequestID(c)
		cv := *models.ExampleCV()

		yearsSinceWork := 3

		dummyData := DataSendToHook{
			MatchedProfiles: []match.FoundMatch{{
				Matches: models.Match{
					M:                 db.NewM(),
					RequestID:         requestID,
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
			CV:        cv,
			KeyID:     mock.Key1.ID,
			KeyName:   "example.com",
			RequestID: requestID,
		}

		dummyDataAsJSON, err := json.Marshal(dummyData)
		if err != nil {
			return err
		}

		err = ctx.GetOnMatchHook(c).Call(bytes.NewReader(dummyDataAsJSON))
		if err != nil {
			return err
		}

		return c.JSON(ExplainDataSendToHook{DataSendToHook: dummyData})
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
			dbConn := ctx.GetDbConn(c)
			hook := models.OnMatchHook{}
			query := bson.M{"_id": hookID}
			args := db.FindOptions{NoDefaultFilters: true}
			err = dbConn.FindOne(&hook, query, args)
			if err != nil {
				return err
			}

			c.SetUserContext(
				ctx.SetOnMatchHook(
					c.UserContext(),
					&hook,
				),
			)

			return c.Next()
		},
	}
}
