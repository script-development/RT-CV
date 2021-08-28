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

func routeGetCvSchema(c *fiber.Ctx) error {
	resSchema, err := schema.From(models.CV{}, "/api/v1/schema/cv")
	if err != nil {
		return err
	}
	return c.JSON(resSchema)
}

func routeGetOpenAPISchema(r *routeBuilder.Router) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		// TODO list the routes in the response
		r.Routes()
		return c.JSON(IMap{
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
			"servers": []IMap{
				{
					"url":         origin,
					"description": "The current server",
				},
			},
			"version": "1.0.0",
		})
	}
}
