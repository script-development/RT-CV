package controller

// This file contains json schemas for certain types
// These are used by the dashboard to to validate user input data
//
// To be specific the monaco editor (the code editor we use in the dashboard) uses json schemas to validate the user input

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/helpers/schema"
	"github.com/script-development/RT-CV/models"
)

var routeGetCvSchema = routeBuilder.R{
	Description: "returns cv as a json schema, " +
		"currently used in the dashboard for the try matcher page to give intelligent suggestions",
	Res: schema.Property{},
	Fn: func(c *fiber.Ctx) error {
		defs := map[string]schema.Property{}
		resSchema, err := schema.From(
			models.CV{},
			"#/$defs/",
			func(key string, value schema.Property) {
				defs[key] = value
			},
			func(key string) bool {
				_, ok := defs[key]
				return ok
			},
			&schema.WithMeta{SchemaID: "/api/v1/schema/cv"},
		)
		if err != nil {
			return err
		}
		resSchema.Defs = defs
		return c.JSON(resSchema)
	},
}

var routeGetOpenAPISchemaCache IMap

func routeGetOpenAPISchema(r *routeBuilder.Router) routeBuilder.R {
	// TODO we use a lot of IMap in here, we should use a typed struct
	return routeBuilder.R{
		Description: "returns openapi as a json schema",
		Res:         IMap{},
		Fn: func(c *fiber.Ctx) error {
			origin := c.Get("Origin")
			if origin == "" {
				// Some web browsers don't have the Origin header so we also check for the host header
				host := c.Get("Host")
				if host == "" {
					host = "localhost:4000"
				}

				toCheck := []string{
					"localhost",
					"127.0.0.1",
					"::1",
				}
				isHTTPS := true
				for _, entry := range toCheck {
					if entry == host || strings.HasPrefix(host, entry+":") {
						isHTTPS = false
						break
					}
				}
				// Do not check c.Secure() as it doesn't work with a proxy
				if isHTTPS {
					origin = "https://" + host
				} else {
					origin = "http://" + host
				}
			}

			// Check if we have a cached version of the response as it's
			if routeGetOpenAPISchemaCache != nil {
				// Replace the only variable that changes form server to server
				routeGetOpenAPISchemaCache["servers"] = []IMap{{"url": origin}}

				return c.JSON(routeGetOpenAPISchemaCache)
			}

			type pathMethods struct {
				Get        IMap   `json:"get,omitempty"`
				Post       IMap   `json:"post,omitempty"`
				Patch      IMap   `json:"patch,omitempty"`
				Put        IMap   `json:"put,omitempty"`
				Delete     IMap   `json:"delete,omitempty"`
				Parameters []IMap `json:"parameters,omitempty"`
			}

			errRes := IMap{
				"description": "unexpected error",
				"content": IMap{
					"application/json": IMap{
						"schema": IMap{
							"$ref": "#/components/schemas/Error",
						},
					},
				},
			}

			componentsSchema := IMap{
				"Error": schema.Property{
					Type:     schema.PropertyTypeObject,
					Required: []string{"error"},
					Properties: map[string]schema.Property{
						"error": {
							Type: schema.PropertyTypeString,
						},
					},
				},
			}

			paths := map[string]pathMethods{}

			for _, route := range r.Routes() {
				contentInfo := IMap{}
				if route.Info.Res != nil {
					schemaValue, err := schema.From(
						route.Info.Res,
						"#/components/schemas/",
						func(key string, value schema.Property) {
							componentsSchema[key] = value
						},
						func(key string) bool {
							_, ok := componentsSchema[key]
							return ok
						},
						nil,
					)
					if err != nil {
						return err
					}
					contentInfo["schema"] = schemaValue
				}

				okRes := IMap{
					"description": "response",
					"content": IMap{
						route.ResponseContentType.String(): contentInfo,
					},
				}

				routeInfo := IMap{
					"summary": strings.TrimPrefix(route.OpenAPIPath, "/api/v1"),
					"responses": IMap{
						"200":     okRes,
						"default": errRes,
					},
				}

				if route.Info.Description != "" {
					routeInfo["description"] = route.Info.Description
				}

				path := paths[route.OpenAPIPath]
				switch route.Method {
				case routeBuilder.Get:
					path.Get = routeInfo
				case routeBuilder.Post:
					path.Post = routeInfo
				case routeBuilder.Patch:
					path.Patch = routeInfo
				case routeBuilder.Put:
					path.Put = routeInfo
				case routeBuilder.Delete:
					path.Delete = routeInfo
				}

			paramsLoop:
				for _, param := range route.Params {
					for _, p := range path.Parameters {
						if p["name"] == param {
							continue paramsLoop
						}
					}
					path.Parameters = append(path.Parameters, IMap{
						"name":     param,
						"in":       "path",
						"required": true,
						"schema": IMap{
							"type": "string",
						},
					})
				}
				paths[route.OpenAPIPath] = path
			}

			res := IMap{
				"openapi": "3.0.3",
				"info": IMap{
					"version":        "1.0.0",
					"title":          "RT-CV",
					"description":    "Real time curriculum vitae matcher",
					"termsOfService": "https://github.com/script-development/RT-CV/blob/main/LICENSE",
					"contact": IMap{
						"name": "API Support",
						"url":  "https://github.com/script-development/RT-CV/issues/new",
					},
					"license": IMap{
						"name": "MIT",
						"url":  "https://github.com/script-development/RT-CV/blob/main/LICENSE",
					},
				},
				"servers":    []IMap{{"url": origin}},
				"paths":      paths,
				"components": IMap{"schemas": componentsSchema},
			}

			// cache the response so we re-use it later on
			routeGetOpenAPISchemaCache = res

			return c.JSON(res)
		},
	}
}
