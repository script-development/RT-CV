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
	Description: "returns cv as a json schema",
	Res:         schema.Property{},
	Fn: func(c *fiber.Ctx) error {
		resSchema, err := schema.From(models.CV{}, "/api/v1/schema/cv")
		if err != nil {
			return err
		}
		return c.JSON(resSchema)
	},
}

func routeGetOpenAPISchema(r *routeBuilder.Router) routeBuilder.R {
	return routeBuilder.R{
		Description: "returns openapi as a json schema",
		Res:         IMap{},
		Fn: func(c *fiber.Ctx) error {
			origin := c.Get("Origin")
			if origin == "" {
				// Some web browsers don't have the Origin header so we also check for the host header
				host := c.Get("Host")
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

			type pathMethods struct {
				Get        IMap `json:"get,omitempty"`
				Post       IMap `json:"post,omitempty"`
				Patch      IMap `json:"patch,omitempty"`
				Put        IMap `json:"put,omitempty"`
				Delete     IMap `json:"delete,omitempty"`
				Parameters IMap `json:"parameters,omitempty"`
			}

			paths := map[string]pathMethods{}
			for _, route := range r.Routes() {
				routeInfo := IMap{
					"responses": IMap{
						"200": IMap{
							"description": "response",
							"content": IMap{
								route.ResponseContentType.String(): IMap{},
							},
						},
						"default": IMap{
							"description": "unexpected error",
							"content": IMap{
								route.ResponseContentType.String(): IMap{},
							},
						},
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

				if len(route.Params) > 0 {
					if path.Parameters == nil {
						path.Parameters = IMap{}
					}
					for _, param := range route.Params {
						path.Parameters[param] = IMap{
							"name":     param,
							"in":       "query",
							"required": true,
						}
					}
				}
				paths[route.OpenAPIPath] = path
			}

			return c.JSON(IMap{
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
				"servers": []IMap{{"url": origin}},
				"paths":   paths,
			})
		},
	}
}
