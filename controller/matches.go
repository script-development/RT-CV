package controller

import (
	"errors"
	"strconv"
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

		now := time.Now()
		if to.Before(now) {
			// We cannot do things in the past, we're fine with caching

			// 1 second * 60 = minute * 60 = hour * 24 = day * 14 = 2 weeks
			maxAge := 1 * 60 * 60 * 24 * 14
			c.Response().Header.Set("Cache-Control", "private")
			c.Response().Header.Add("Cache-Control", "max-age="+strconv.Itoa(maxAge))

			c.Response().Header.SetLastModified(to)
		}

		query := bson.M{
			"when": bson.M{
				"$gt": from,
				"$lt": to,
			},
		}

		if c.Params(`profile`) != "" {
			query["profileId"] = ctx.GetProfile(c).ID
		}

		matches := []models.Match{}
		err = db.Find(&models.Match{}, &matches, query)
		if err != nil {
			return err
		}

		return c.JSON(matches)
	},
}
