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

	Routes(app, db, testingServerSeed)
	return &testingRouter{
		t:        t,
		fiber:    app,
		db:       db,
		accessor: auth.NewAccessorHelper(mock.Key1.ID, "abc", string(random.StringBytes(32)), testingServerSeed),
	}
}

type TestReqOpts struct {
	NoAuth bool
	Body   []byte
}

func (r *testingRouter) MakeRequest(method Method, route string, opts TestReqOpts) (res *http.Response, resBody []byte) {
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
		method Method
		route  string
	}{
		{
			"scraper",
			Post,
			"/api/v1/scraper/scanCV",
		},
		{
			"control",
			Get,
			"/api/v1/control/reloadProfiles",
		},
	}

	app := newTestingRouter(t)

	for _, route := range routes {
		route := route

		t.Run(route.name, func(t *testing.T) {
			res, body := app.MakeRequest(route.method, route.route, TestReqOpts{
				NoAuth: true,
			})

			// 401 = Unauthorized
			Equal(t, 401, res.StatusCode)
			Equal(t, `{"error":"missing authorization header of type Basic"}`, string(body))
		})
	}
}
