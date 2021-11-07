package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

var errAuthMissingRoles = errors.New("you do not have auth roles required to access this route")

func requiresAuth(requiredRoles models.APIKeyRole) routeBuilder.M {
	tags := []routeBuilder.Tag{
		{
			Name:        "Auth all route",
			Description: "route that requires authentication",
		},
	}

	requiredRolesList := requiredRoles.ConvertToAPIRoles()
	for _, role := range requiredRolesList {
		roleStr := strconv.FormatUint(uint64(role.Role), 10)
		tags = append(tags, routeBuilder.Tag{
			Name: "Auth Role " + role.Slug,
			Description: fmt.Sprintf(
				"route required authentication id %s, description: %s",
				roleStr,
				role.Description,
			),
		})
	}

	return routeBuilder.M{
		Tags: tags,
		Fn: func(c *fiber.Ctx) error {
			key := ctx.GetKey(c)
			// Check if the auth header is already checked earlier in the request
			// If true we only have to check if the roles match
			if key != nil {
				if requiredRoles != 0 && !key.Roles.ContainsSome(requiredRoles) {
					return ErrorRes(c, 401, errAuthMissingRoles)
				}
				return c.Next()
			}

			// Get values from context
			authService := ctx.GetAuth(c)
			logger := ctx.GetLogger(c)

			// Check auth header
			authorizationValue := c.Get("Authorization")
			if len(authorizationValue) == 0 {
				// Used by the websocket routes
				authorizationValue = c.Params("authorization")
				if len(authorizationValue) != 0 {
					authorizationValue = "Basic " + authorizationValue
				}
			}
			if len(authorizationValue) == 0 {
				// NOTE: there seems to be a bug with fiber it seems where if try to access a non existing route or send an invalid url
				// c.Get("Authorization") returns an empty string no matter the value of the header send
				// This might cause some confusuion as you'll receive a auth.ErrNoAuthheader error over a 404 error
				return ErrorRes(c, fiber.StatusBadRequest, auth.ErrNoAuthHeader)
			}

			key, err := authService.Valid(authorizationValue)
			if err != nil {
				return ErrorRes(c, fiber.StatusUnauthorized, err)
			}

			// Check required roles matches
			if requiredRoles != 0 && !key.Roles.ContainsSome(requiredRoles) {
				return ErrorRes(c, fiber.StatusForbidden, errAuthMissingRoles)
			}

			*logger = *logger.WithFields(log.Fields{
				"api_key_id": key.ID.Hex(),
				"domains":    key.Domains,
			})

			c.SetUserContext(
				ctx.SetKey(
					c.UserContext(),
					key,
				),
			)

			return c.Next()
		},
	}
}
