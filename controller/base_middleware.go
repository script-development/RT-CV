package controller

import (
	"context"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/models"
)

// InsertData adds the profiles to every route
func InsertData(dbConn db.Connection) fiber.Handler {
	profiles, err := models.GetProfiles(dbConn)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.WithValue(context.Background(), profilesCtxKey, &profiles)

	keys, err := models.GetApiKeys(dbConn)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx = context.WithValue(ctx, authCtxKey, auth.New(keys))

	ctx = context.WithValue(ctx, dbConnCtxKey, dbConn)

	// Pre define loggerEntity so we only take once memory
	loggerEntity := log.Entry{
		Logger: log.Log.(*log.Logger),
	}

	return func(c *fiber.Ctx) error {
		// reset loggerEntity
		loggerEntity = log.Entry{
			Logger: loggerEntity.Logger,
		}

		c.SetUserContext(context.WithValue(ctx, loggerCtxKey, &loggerEntity))
		return c.Next()
	}
}
