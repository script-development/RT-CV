package routeBuilder

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Tag is meta information for a route
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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
	Tags []Tag
}

// M gives this package more information about a middleware
type M struct {
	Tags []Tag
	Fn   fiber.Handler
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
	tags   []Tag
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

func appendTags(tags, other []Tag) []Tag {
outerLoop:
	for _, otherTag := range other {
		// Firstly lets check for tag duplicates
		for _, tag := range tags {
			if tag.Name == otherTag.Name {
				if tag.Description != otherTag.Description {
					msg := fmt.Sprintf(
						"found 2 tags with the same name but diffrent description, "+
							"tagname: %s, description 1: %s, description 2: %s",
						tag.Name,
						tag.Description,
						otherTag.Description,
					)
					panic(msg)
				}
				continue outerLoop
			}
		}

		tags = append(tags, otherTag)
	}
	return tags
}

func (r *Router) appendTags(middlewares []M) []Tag {
	tags := r.tags
	for _, middleware := range middlewares {
		tags = appendTags(tags, middleware.Tags)
	}
	return tags
}

func (r *Router) newRoute(prefix string, method Method, info R, middlewares []M) {
	fiberPath := r.appendPrefix(prefix)
	parsedPath := parseFiberPath(fiberPath)

	// Append all middleware tags to Info.Tags
	info.Tags = appendTags(r.appendTags(middlewares), info.Tags)

	r.base.routes = append(r.base.routes, Route{
		FiberPath:           r.appendPrefix(prefix),
		OpenAPIPath:         parsedPath.AsOpenAPIPath,
		Params:              parsedPath.Params,
		Method:              method,
		ResponseContentType: JSON,
		Info:                info,
	})
}

func getHandlers(middleware []M, route *R) []func(*fiber.Ctx) error {
	handlers := make([]func(*fiber.Ctx) error, len(middleware))
	for idx, middleware := range middleware {
		handlers[idx] = middleware.Fn
	}
	if route != nil {
		handlers = append(handlers, route.Fn)
	}
	return handlers
}

// Group prefixes the routes within the group with a route and adds a middleware to them if specified
func (r *Router) Group(prefix string, group func(*Router), middlewares ...M) {
	group(&Router{
		tags:   r.appendTags(middlewares),
		prefix: r.appendPrefix(prefix),
		fiber:  r.fiber.Group(prefix, getHandlers(middlewares, nil)...),
		base:   r.base,
	})
}

// Get defines a get route with information about the route
func (r *Router) Get(prefix string, routeDefinition R, middlewares ...M) {
	routeDefinition.check()
	r.newRoute(prefix, Get, routeDefinition, middlewares)
	r.fiber.Get(prefix, getHandlers(middlewares, &routeDefinition)...)
}

// Post defines a POST route
func (r *Router) Post(prefix string, routeDefinition R, middlewares ...M) {
	routeDefinition.check()
	r.newRoute(prefix, Post, routeDefinition, middlewares)
	r.fiber.Post(prefix, getHandlers(middlewares, &routeDefinition)...)
}

// Put defines a PUT route
func (r *Router) Put(prefix string, routeDefinition R, middlewares ...M) {
	routeDefinition.check()
	r.newRoute(prefix, Put, routeDefinition, middlewares)
	r.fiber.Put(prefix, getHandlers(middlewares, &routeDefinition)...)
}

// Patch defines a PATCH route
func (r *Router) Patch(prefix string, routeDefinition R, middlewares ...M) {
	routeDefinition.check()
	r.newRoute(prefix, Patch, routeDefinition, middlewares)
	r.fiber.Patch(prefix, getHandlers(middlewares, &routeDefinition)...)
}

// Delete defines a DELETE route
func (r *Router) Delete(prefix string, routeDefinition R, middlewares ...M) {
	routeDefinition.check()
	r.newRoute(prefix, Delete, routeDefinition, middlewares)
	r.fiber.Delete(prefix, getHandlers(middlewares, &routeDefinition)...)
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
