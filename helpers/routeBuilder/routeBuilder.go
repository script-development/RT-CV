package routeBuilder

import (
	"github.com/gofiber/fiber/v2"
)

// R gives this package more information about a route
// From the fiber handler (fn(c *fiber.Ctx) error) we cannot know the expected input and output data
type R struct {
	// Required
	Description string
	// Res gives the routes
	Res interface{}
	Fn  func(c *fiber.Ctx) error

	// Optional
	// Body can be set to hint the route builder what kind of body this request expects
	Body interface{}
}

func (r R) check() {
	if r.Description == "" {
		panic("routeBuilder.R description cannot be empty")
	}
	if r.Res == nil {
		panic("routeBuilder.Res must be defined")
	}
	if r.Fn == nil {
		panic("routeBuilder.Fn must be defined")
	}
}

// BaseBuilder is the core builder in here all the routes and middlewares are rememberd
type BaseBuilder struct {
	fiber  *fiber.App
	routes []Route
}

// Method represends a http method
type Method uint8

const (
	// Get is a http method
	Get Method = iota
	// Post is a http method
	Post
	// Patch is a http method
	Patch
	// Put is a http method
	Put
	// Delete is a http method
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

// ContentType represends a content type
type ContentType uint8

const (
	// JSON is a content type
	JSON ContentType = iota
	// HTML is a content type
	HTML
)

func (c ContentType) String() string {
	switch c {
	case JSON:
		return "application/json"
	case HTML:
		return "text/html"
	default:
		return "text/plain"
	}
}

// Route constains information about a route
type Route struct {
	FiberPath           string
	OpenAPIPath         string
	Params              []string
	Kind                string
	Method              Method
	ResponseContentType ContentType
	Info                R
}

// New creates a instance of Builder
func New(app *fiber.App) *Router {
	return &Router{
		fiber:  app,
		prefix: "",
		base: &BaseBuilder{
			fiber:  app,
			routes: []Route{},
		},
	}
}

// Router can be used to define routes and middlwares
type Router struct {
	prefix string
	fiber  fiber.Router
	base   *BaseBuilder
}

func (r *Router) appendPrefix(add string) string {
	if len(add) > 0 && add[len(add)-1] == '/' {
		add = add[:len(add)-1]
	}
	if add == "" {
		return r.prefix
	}
	if add[0] != '/' {
		add = "/" + add
	}
	return r.prefix + add
}

func (r *Router) newRoute(prefix string, method Method, info R) {
	fiberPath := r.appendPrefix(prefix)
	parsedPath := parseFiberPath(fiberPath)

	r.base.routes = append(r.base.routes, Route{
		FiberPath:           r.appendPrefix(prefix),
		OpenAPIPath:         parsedPath.AsOpenAPIPath,
		Params:              parsedPath.Params,
		Method:              method,
		ResponseContentType: JSON,
		Info:                info,
	})
}

// Group prefixes the routes within the group with a route and adds a middleware to them if specified
func (r *Router) Group(prefix string, group func(*Router), middlewares ...func(*fiber.Ctx) error) {
	group(&Router{
		prefix: r.appendPrefix(prefix),
		fiber:  r.fiber.Group(prefix, middlewares...),
		base:   r.base,
	})
}

// Get defines a get route with information about the route
func (r *Router) Get(prefix string, routeDefinition R, middlewares ...func(*fiber.Ctx) error) {
	routeDefinition.check()
	r.newRoute(prefix, Get, routeDefinition)
	r.fiber.Get(prefix, append(middlewares, routeDefinition.Fn)...)
}

// Post defines a POST route
func (r *Router) Post(prefix string, routeDefinition R, handlers ...func(*fiber.Ctx) error) {
	routeDefinition.check()
	r.newRoute(prefix, Post, routeDefinition)
	r.fiber.Post(prefix, append(handlers, routeDefinition.Fn)...)
}

// Put defines a PUT route
func (r *Router) Put(prefix string, routeDefinition R, handlers ...func(*fiber.Ctx) error) {
	routeDefinition.check()
	r.newRoute(prefix, Put, routeDefinition)
	r.fiber.Put(prefix, append(handlers, routeDefinition.Fn)...)
}

// Patch defines a PATCH route
func (r *Router) Patch(prefix string, routeDefinition R, handlers ...func(*fiber.Ctx) error) {
	routeDefinition.check()
	r.newRoute(prefix, Patch, routeDefinition)
	r.fiber.Patch(prefix, append(handlers, routeDefinition.Fn)...)
}

// Delete defines a DELETE route
func (r *Router) Delete(prefix string, handlers ...func(*fiber.Ctx) error) {
	r.newRoute(prefix, Delete, R{})
	r.fiber.Delete(prefix, handlers...)
}

// Static defines a static file path
// We also don't store the static resources as they are not really important to the api users
func (r *Router) Static(prefix, root string, options ...fiber.Static) {
	r.fiber.Static(prefix, root, options...)
}

// Routes returns all routes
func (r *Router) Routes() []Route {
	return r.base.routes
}
