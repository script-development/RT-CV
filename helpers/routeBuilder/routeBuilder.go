package routeBuilder

import "github.com/gofiber/fiber/v2"

// BaseBuilder is the core builder in here all the routes and middlewares are rememberd
type BaseBuilder struct {
	fiber  *fiber.App
	routes []Route
}

// Route constains information about a route
type Route struct {
	path   string
	method string // GET, POST, PATCH, DELETE
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

// Group prefixes the routes within the group with a route and adds a middlware to them if spesified
func (r *Router) Group(prefix string, group func(*Router), middlewares ...func(*fiber.Ctx) error) {
	group(&Router{
		prefix: r.appendPrefix(prefix),
		fiber:  r.fiber.Group(prefix, middlewares...),
		base:   r.base,
	})
}

// Get defines a GET route
func (r *Router) Get(prefix string, handlers ...func(*fiber.Ctx) error) {
	r.base.routes = append(r.base.routes, Route{
		path:   r.appendPrefix(prefix),
		method: "GET",
	})
	r.fiber.Get(prefix, handlers...)
}

// Post defines a POST route
func (r *Router) Post(prefix string, handlers ...func(*fiber.Ctx) error) {
	r.base.routes = append(r.base.routes, Route{
		path:   r.appendPrefix(prefix),
		method: "POST",
	})
	r.fiber.Post(prefix, handlers...)
}

// Put defines a PUT route
func (r *Router) Put(prefix string, handlers ...func(*fiber.Ctx) error) {
	r.base.routes = append(r.base.routes, Route{
		path:   r.appendPrefix(prefix),
		method: "PUT",
	})
	r.fiber.Put(prefix, handlers...)
}

// Patch defines a PATCH route
func (r *Router) Patch(prefix string, handlers ...func(*fiber.Ctx) error) {
	r.base.routes = append(r.base.routes, Route{
		path:   r.appendPrefix(prefix),
		method: "PATCH",
	})
	r.fiber.Patch(prefix, handlers...)
}

// Delete defines a DELETE route
func (r *Router) Delete(prefix string, handlers ...func(*fiber.Ctx) error) {
	r.base.routes = append(r.base.routes, Route{
		path:   r.appendPrefix(prefix),
		method: "DELETE",
	})
	r.fiber.Delete(prefix, handlers...)
}

// Static defines a static file path
func (r *Router) Static(prefix, root string, options ...fiber.Static) {
	r.base.routes = append(r.base.routes, Route{
		path:   r.appendPrefix(prefix),
		method: "GET",
	})
	r.fiber.Static(prefix, root, options...)
}

// Routes returns all routes
func (r *Router) Routes() []Route {
	return r.base.routes
}
