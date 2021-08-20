package controller

import (
	"context"
	"errors"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/models"
)

func requiresAuth(requiredRoles models.ApiKeyRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get values from context
		auth := GetAuth(c)
		logger := GetLogger(c)

		// Check auth header
		authorizationHeader := c.Get("Authorization")
		key, salt, err := auth.Authenticate([]byte(authorizationHeader))
		if err != nil {
			return ErrorRes(c, 401, err)
		}

		// Check required roles matches
		if !key.Roles.ContainsSome(requiredRoles) {
			return ErrorRes(c, 401, errors.New("you do not have the permissions to access this route"))
		}

		*logger = *logger.WithFields(log.Fields{
			"apiKey":     key.Key,
			"apiKeySalt": string(salt),
			"domains":    key.Domains,
		})

		c.SetUserContext(
			context.WithValue(
				c.UserContext(),
				keyCtxKey,
				key,
			),
		)

		return c.Next()
	}
}
