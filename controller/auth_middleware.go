package controller

import (
	"context"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

func requiresAuth(requiredRoles models.ApiKeyRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get values from context
		auth := ctx.Value(AuthCtxKey).(*auth.Auth)
		logger := ctx.Value(LoggerCtxKey).(*log.Entry)

		// Check auth header
		authorizationHeader := []byte(c.Get("Authorization"))
		key, salt, err := auth.Authenticate(authorizationHeader)
		if err != nil {
			return c.Status(401).SendString(err.Error())
		}

		// Check required roles matches
		if !key.Roles.ContainsSome(requiredRoles) {
			return c.Status(401).SendString("you do not have the permissions to access this route")
		}

		ctx = context.WithValue(ctx, KeyCtxKey, key)

		*logger = *logger.WithFields(log.Fields{
			"apiKey":     key.Key,
			"apiKeySalt": string(salt),
			"domains":    key.Domains,
		})

		c.SetUserContext(ctx)
		return c.Next()
	}
}
