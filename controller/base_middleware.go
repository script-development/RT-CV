package controller

import (
	"context"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/auth"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertData adds the profiles to every route
func InsertData(dbConn db.Connection) routeBuilder.M {
	requestContext := ctx.SetDbConn(context.Background(), dbConn)

	requestContext = ctx.ResetMatcherProfilesCache(requestContext)

	requestContext = ctx.SetAuth(requestContext, auth.NewHelper(dbConn))

	// We set this to nil so we can later run ctx.GetKey without panicing if the key is not yet set
	requestContext = ctx.SetKey(requestContext, nil)

	// Pre define loggerEntity so we only take once memory
	loggerEntity := &log.Entry{
		Logger: log.Log.(*log.Logger),
	}

	return routeBuilder.M{
		Fn: func(c *fiber.Ctx) error {
			requestID := primitive.NewObjectID()
			c.Response().Header.Add("X-Request-ID", requestID.Hex())

			loggerEntity = log.WithField("request_id", requestID.Hex())

			c.SetUserContext(
				ctx.SetRequestID(
					ctx.SetLogger(
						requestContext,
						loggerEntity,
					),
					requestID,
				),
			)
			return c.Next()
		},
	}
}
