package controller

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
)

var routeGetMatchesPeriod = routeBuilder.R{
	Description: `get all matches made within a certain period.
the from and to param should be in the RFC 3339 format
_RFC 3339 is basically an extension to iso 8601_`,
	Res: []models.Match{},
	Fn: func(c *fiber.Ctx) error {
		db := ctx.GetDbConn(c)

		fromParam, toParam := c.Params("from"), c.Params("to")
		from, err := time.Parse(time.RFC3339, fromParam)
		if err != nil {
			return errors.New("unable \"from\" parse from param as RFC 3339")
		}
		to, err := time.Parse(time.RFC3339, toParam)
		if err != nil {
			return errors.New("unable \"to\" parse from param as RFC 3339")
		}

		query := bson.M{
			"when": bson.M{
				"$gt": from,
				"$lt": to,
			},
		}

		matches := []models.Match{}
		err = db.Find(&models.Match{}, &matches, query)
		if err != nil {
			return err
		}

		return c.JSON(matches)
	},
}
