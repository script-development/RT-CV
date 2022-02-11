package controller

// This file contains json schemas for certain types
// These are used by the dashboard to to validate user input data
//
// To be specific the monaco editor (the code editor we use in the dashboard) uses json schemas to validate the user input

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mjarkk/jsonschema"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

var routeGetCvSchema = routeBuilder.R{
	Description: "returns cv as a json schema, " +
		"currently used in the dashboard for the /tryMatcher page to give intelligent suggestions",
	Res: jsonschema.Property{},
	Fn: func(c *fiber.Ctx) error {
		defs := map[string]jsonschema.Property{}
		resSchema, err := jsonschema.From(
			models.CV{},
			"#/$defs/",
			func(key string, value jsonschema.Property) {
				defs[key] = value
			},
			func(key string) bool {
				_, ok := defs[key]
				return ok
			},
			&jsonschema.WithMeta{SchemaID: "/api/v1/schema/cv"},
		)
		if err != nil {
			return err
		}
		resSchema.Defs = defs

		return c.JSON(resSchema)
	},
}

var errResponse = routeBuilder.OpenAPIResponse{
	Description: "unexpected error",
	Content: map[string]routeBuilder.OpenAPIMediaType{
		"application/json": {
			Schema: IMap{
				"$ref": "#/components/schemas/Error",
			},
		},
	},
}

var routeGetOpenAPISchemaCache *routeBuilder.OpenAPI

func routeGetOpenAPISchema(r *routeBuilder.Router) routeBuilder.R {
	// FIXME replace IMap with structs
	return routeBuilder.R{
		Description: "Returns the api schema as an openapi schema\n" +
			"This schema is currently used by the /docs page",
		Res: routeBuilder.OpenAPI{},
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
				schemaInCache := *routeGetOpenAPISchemaCache
				schemaInCache.Servers = []routeBuilder.OpenAPIServer{{URL: origin}}
				return c.JSON(schemaInCache)
			}

			componentsSchema := IMap{
				"Error": jsonschema.Property{
					Type:     jsonschema.PropertyTypeObject,
					Required: []string{"error"},
					Properties: map[string]jsonschema.Property{
						"error": {
							Type: jsonschema.PropertyTypeString,
						},
					},
				},
			}

			paths := map[string]routeBuilder.OpenAPIPathItem{}

			allTags := []routeBuilder.Tag{}

			for _, route := range r.Routes() {
				responsesMap := map[string]routeBuilder.OpenAPIResponse{
					"error": errResponse,
				}

				// Create the response value
				if route.Info.Res != nil {
					schemaValue, err := jsonschema.From(
						route.Info.Res,
						"#/components/schemas/",
						func(key string, value jsonschema.Property) {
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

					responsesMap["200"] = routeBuilder.OpenAPIResponse{
						Description: "response",
						Content: map[string]routeBuilder.OpenAPIMediaType{
							route.ResponseContentType.String(): {Schema: schemaValue},
						},
					}
				} else if route.Info.ResMap != nil {
					for key, value := range route.Info.ResMap {
						content := routeBuilder.OpenAPIMediaType{}
						if value != nil {
							schemaValue, err := jsonschema.From(
								value,
								"#/components/schemas/",
								func(key string, value jsonschema.Property) {
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
							content.Schema = schemaValue
						}

						responsesMap[key] = routeBuilder.OpenAPIResponse{
							Description: "response",
							Content: map[string]routeBuilder.OpenAPIMediaType{
								route.ResponseContentType.String(): content,
							},
						}
					}
				} else if route.Info.CustomResponse != nil {
					responsesMap["200"] = *route.Info.CustomResponse
				} else {
					responsesMap["200"] = routeBuilder.OpenAPIResponse{
						Description: "response",
						Content: map[string]routeBuilder.OpenAPIMediaType{
							route.ResponseContentType.String(): {},
						},
					}
				}

				// Create the actual information about this route's method
				routeInfo := routeBuilder.OpenAPIOperation{
					Summary:   strings.TrimPrefix(route.OpenAPIPath, "/api/v1"),
					Responses: responsesMap,
				}
				if route.Info.Description != "" {
					routeInfo.Description = route.Info.Description
				}

				if len(route.Info.Tags) > 0 {
					tagsList := make([]string, len(route.Info.Tags))
					for idx, tag := range route.Info.Tags {
						tagsList[idx] = tag.Name

						alreadyDefined := false
						for _, tagFromAllTags := range allTags {
							if tagFromAllTags.Name == tag.Name {
								alreadyDefined = true
								break
							}
						}
						if !alreadyDefined {
							allTags = append(allTags, tag)
						}
					}
					routeInfo.Tags = tagsList
				}

				// Create the request body expected value
				if route.Info.Body != nil {
					schemaValue, err := jsonschema.From(
						route.Info.Body,
						"#/components/schemas/",
						func(key string, value jsonschema.Property) {
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
					routeInfo.RequestBody = &routeBuilder.OpenAPIRequestBody{
						Description: "request data",
						Content: map[string]routeBuilder.OpenAPIMediaType{
							"application/json": {Schema: schemaValue},
						},
					}
				}

				// Insert the above created routeInfo
				path := paths[route.OpenAPIPath]
				switch route.Method {
				case routeBuilder.Get:
					path.Get = &routeInfo
				case routeBuilder.Post:
					path.Post = &routeInfo
				case routeBuilder.Patch:
					path.Patch = &routeInfo
				case routeBuilder.Put:
					path.Put = &routeInfo
				case routeBuilder.Delete:
					path.Delete = &routeInfo
				}

			paramsLoop:
				for _, param := range route.Params {
					// Insert the url params
					for _, p := range path.Parameters {
						if p.Name == param {
							continue paramsLoop
						}
					}
					path.Parameters = append(path.Parameters, routeBuilder.OpenAPIParameter{
						Name:     param,
						In:       "path",
						Required: true,
						Schema: IMap{
							"type": "string",
						},
					})
				}
				paths[route.OpenAPIPath] = path
			}

			res := routeBuilder.OpenAPI{
				OpenAPI: "3.0.3",
				Info: &routeBuilder.OpenAPIInfo{
					Version:        "1.0.0",
					Title:          "RT-CV",
					Description:    schemaDescription,
					TermsOfService: "https://github.com/script-development/RT-CV/blob/main/LICENSE",
					Contact: routeBuilder.OpenAPIContact{
						Name: "API Support",
						URL:  "https://github.com/script-development/RT-CV/issues/new",
					},
					License: routeBuilder.OpenAPILicense{
						Name: "MIT",
						URL:  "https://github.com/script-development/RT-CV/blob/main/LICENSE",
					},
				},
				Servers:    []routeBuilder.OpenAPIServer{{URL: origin}},
				Paths:      paths,
				Components: &routeBuilder.OpenAPIComponents{Schemas: componentsSchema},
				Tags:       allTags,
			}

			// cache the response so we re-use it later on
			routeGetOpenAPISchemaCache = &res

			return c.JSON(res)
		},
	}
}
