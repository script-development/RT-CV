package controller

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

var scannedReferenceNrs = routeBuilder.R{
	Description: "get a list of all earlier scraped reference numbers",
	Res:         []string{},
	Fn: func(c *fiber.Ctx) error {
		hours := c.Params("hours")
		days := c.Params("days")
		weeks := c.Params("weeks")

		dbConn := ctx.GetDbConn(c)
		key := ctx.GetKey(c)

		matches := []models.Match{}
		var err error
		now := time.Now()

		switch {
		case hours != "":
			hoursInt, err := strconv.Atoi(hours)
			if err != nil {
				return errors.New("hours argument is not a valid number")
			}
			if hoursInt <= 0 {
				return errors.New("hours argument must be greater than 0")
			}
			matches, err = models.GetMatchesSince(dbConn, now.Add(-(time.Hour * time.Duration(hoursInt))), &key.ID)
		case days != "":
			daysInt, err := strconv.Atoi(days)
			if err != nil {
				return errors.New("days argument is not a valid number")
			}
			if daysInt <= 0 {
				return errors.New("days argument must be greater than 0")
			}
			matches, err = models.GetMatchesSince(dbConn, now.AddDate(0, 0, -daysInt), &key.ID)
		case weeks != "":
			weeksInt, err := strconv.Atoi(weeks)
			if err != nil {
				return errors.New("weeks argument is not a valid number")
			}
			if weeksInt <= 0 {
				return errors.New("weeks argument must be greater than 0")
			}
			matches, err = models.GetMatchesSince(dbConn, now.AddDate(0, 0, -(7*weeksInt)), &key.ID)
		default:
			matches, err = models.GetMatches(dbConn, &key.ID)
		}
		if err != nil {
			return err
		}

		res := []string{}
	outerLoop:
		for idx, ref := range matches {
			refNr := ref.ReferenceNr
			if idx == 0 {
				res = append(res, refNr)
				continue
			}

			// Check for duplicated reference numbers
			for idx := len(res) - 1; idx >= 0; idx-- {
				if res[idx] == refNr {
					// This reference nr is already in the respnose list
					continue outerLoop
				}
			}

			res = append(res, refNr)
		}

		return c.JSON(res)
	},
}
