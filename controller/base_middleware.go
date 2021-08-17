package controller

import (
	"context"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

// InsertData adds the profiles to every route
func InsertData() fiber.Handler {
	profiles, err := models.GetProfiles()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.WithValue(context.Background(), ProfilesCtxKey, &profiles)

	keys, err := models.GetApiKeys()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx = context.WithValue(ctx, AuthCtxKey, auth.New(keys))

	// Pre define loggerEntity so we only take once memory
	loggerEntity := log.Entry{
		Logger: log.Log.(*log.Logger),
	}

	return func(c *fiber.Ctx) error {
		// reset loggerEntity
		loggerEntity = log.Entry{
			Logger: loggerEntity.Logger,
		}

		c.SetUserContext(context.WithValue(ctx, LoggerCtxKey, &loggerEntity))
		return c.Next()
	}
}
