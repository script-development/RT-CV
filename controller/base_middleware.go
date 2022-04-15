package controller

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertData adds the profiles to every route
func InsertData(dbConn db.Connection) fiber.Handler {
	authHelper := auth.NewHelper(dbConn)

	// Pre define loggerEntity so we only take once memory
	loggerEntity := &log.Entry{
		Logger: log.Log.(*log.Logger),
	}

	matcherProfilesCache := &ctx.MatcherProfilesCache{}

	return func(c *fiber.Ctx) error {
		requestID := primitive.NewObjectID()
		c.Response().Header.Add("X-Request-ID", requestID.Hex())

		c.SetUserContext(ctx.Set(c.UserContext(), &ctx.Ctx{
			RequestID:            requestID,
			Auth:                 authHelper,
			Logger:               loggerEntity.WithField("request_id", requestID.Hex()),
			DBConn:               dbConn,
			MatcherProfilesCache: matcherProfilesCache,
		}))

		return c.Next()
	}
}
