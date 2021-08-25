package controller

import (
	"errors"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/models"
)

var errAuthMissingRoles = errors.New("you do not have auth roles required to access this route")

func requiresAuth(requiredRoles models.APIKeyRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		auth := ctx.GetAuth(c)
		logger := ctx.GetLogger(c)

		// Check auth header
		authorizationHeader := c.Get("Authorization")
		key, salt, err := auth.Authenticate([]byte(authorizationHeader))
		if err != nil {
			return ErrorRes(c, 401, err)
		}

		// Check required roles matches
		if requiredRoles != 0 && !key.Roles.ContainsSome(requiredRoles) {
			return ErrorRes(c, 401, errAuthMissingRoles)
		}

		*logger = *logger.WithFields(log.Fields{
			"apiKey":     key.Key,
			"apiKeySalt": string(salt),
			"domains":    key.Domains,
		})

		c.SetUserContext(
			ctx.SetKey(
				c.UserContext(),
				key,
			),
		)

		return c.Next()
	}
}
