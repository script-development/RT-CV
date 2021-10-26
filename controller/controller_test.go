package controller

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db/testingdb"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/helpers/random"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	. "github.com/stretchr/testify/assert"
)

var testingServerSeed = "static-server-seed"

type testingRouter struct {
	t        *testing.T
	fiber    *fiber.App
	db       *testingdb.TestConnection
	accessor *auth.TestAccessor
}

func newTestingRouter(t *testing.T) *testingRouter {
	db := mock.NewMockDB()
	app := fiber.New(fiber.Config{
		ErrorHandler: FiberErrorHandler,
	})

	Routes(app, "TESTING", db, testingServerSeed, true)
	return &testingRouter{
		t:        t,
		fiber:    app,
		db:       db,
		accessor: auth.NewAccessorHelper(mock.Key1.ID, "aaa", string(random.StringBytes(32)), testingServerSeed),
	}
}

type TestReqOpts struct {
	NoAuth bool
	Body   []byte
}

func (r *testingRouter) MakeRequest(method routeBuilder.Method, route string, opts TestReqOpts) (res *http.Response, resBody []byte) {
	var body io.Reader
	if opts.Body != nil {
		body = bytes.NewReader(opts.Body)
	}

	req, err := http.NewRequest(
		method.String(),
		route,
		body,
	)
	NoError(r.t, err)

	if opts.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if !opts.NoAuth {
		req.Header.Set("Authorization", string(r.accessor.Key()))
	}

	res, err = r.fiber.Test(req, -1)
	NoError(r.t, err)

	resBody, err = ioutil.ReadAll(res.Body)
	NoError(r.t, err)

	return res, resBody
}

func TestCannotAccessCriticalRoutesWithoutCredentials(t *testing.T) {
	routes := []struct {
		name   string
		method routeBuilder.Method
		route  string
	}{
		{
			"scraper",
			routeBuilder.Post,
			"/api/v1/scraper/scanCV",
		},
		{
			"control profiles",
			routeBuilder.Get,
			"/api/v1/control/profiles",
		},
		{
			"keys",
			routeBuilder.Get,
			"/api/v1/keys",
		},
	}

	app := newTestingRouter(t)

	for _, route := range routes {
		route := route

		t.Run(route.name, func(t *testing.T) {
			res, body := app.MakeRequest(route.method, route.route, TestReqOpts{
				NoAuth: true,
			})

			Equal(t, 400, res.StatusCode, route.route)
			Equal(t, `{"error":"missing authorization header of type Basic"}`, string(body), route.route)
		})
	}
}
