package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

type Method uint

const (
	Get Method = iota
	Post
	Patch
	Put
	Delete
)

func (m Method) String() string {
	switch m {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Patch:
		return "PATCH"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	default:
		return "GET"
	}
}

func TestCannotAccessCriticalRoutesWithoutCredentials(t *testing.T) {
	models.Testing = true

	scenarios := []struct {
		name   string
		method Method
		route  string
	}{
		{
			"scraper",
			Post,
			"/v1/scraper/scanCV",
		},
		{
			"control",
			Get,
			"/v1/control/reloadProfiles",
		},
	}

	app := fiber.New()
	controller.Routes(app)

	for _, s := range scenarios {
		s := s

		t.Run(s.name, func(t *testing.T) {
			req, err := http.NewRequest(
				s.method.String(),
				s.route,
				nil,
			)
			NoError(t, err)

			res, err := app.Test(req, -1)
			NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			NoError(t, err)
			Equal(t, "missing authorization header of type Basic", string(body))

			// 401 = Unauthorized
			Equal(t, 401, res.StatusCode)
		})
	}
}
