package controller

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db/testingdb"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/helpers/random"
	"github.com/script-development/RT-CV/mock"
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

var testingServerSeed = []byte("static-server-seed")

type testingRouter struct {
	fiber    *fiber.App
	db       *testingdb.TestConnection
	accessor *auth.TestAccessor
}

func newTestingRouter() *testingRouter {
	db := mock.NewMockDB()
	app := fiber.New(fiber.Config{
		ErrorHandler: FiberErrorHandler,
	})

	Routes(app, db, testingServerSeed)
	return &testingRouter{
		fiber:    app,
		db:       db,
		accessor: auth.NewAccessorHelper(mock.Key1.ID, "abc", string(random.StringBytes(32)), testingServerSeed),
	}
}

type TestReqOpts struct {
	NoAuth bool
}

func (r *testingRouter) MakeRequest(t *testing.T, method Method, route string, opts TestReqOpts) *http.Response {
	req, err := http.NewRequest(
		Get.String(),
		route,
		nil,
	)
	NoError(t, err)

	if !opts.NoAuth {
		req.Header.Add("Authorization", string(r.accessor.Key()))
	}

	res, err := r.fiber.Test(req, -1)
	NoError(t, err)

	return res
}

func TestCannotAccessCriticalRoutesWithoutCredentials(t *testing.T) {
	routes := []struct {
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

	app := newTestingRouter()

	for _, route := range routes {
		route := route

		t.Run(route.name, func(t *testing.T) {
			res := app.MakeRequest(t, route.method, route.route, TestReqOpts{
				NoAuth: true,
			})

			// 401 = Unauthorized
			Equal(t, 401, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			NoError(t, err)
			Equal(t, `{"error":"missing authorization header of type Basic"}`, string(body))
		})
	}
}
