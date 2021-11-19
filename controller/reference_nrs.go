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

var scannedReferenceNrs = routeBuilder.R{
	Description: "get a list of all earlier scraped reference numbers",
	Res:         []string{},
	Fn: func(c *fiber.Ctx) error {
		hours := c.Params("hours")
		days := c.Params("days")
		weeks := c.Params("weeks")

		dbConn := ctx.GetDbConn(c)
		key := ctx.GetKey(c)

		query := bson.M{"keyId": key.ID}

		switch {
		case hours != "":
			hoursInt, err := strconv.Atoi(hours)
			if err != nil {
				return errors.New("hours argument is not a valid number")
			}
			if hoursInt <= 0 {
				return errors.New("hours argument must be greater than 0")
			}
			query["insertionDate"] = bson.M{"$gt": time.Now().Add(-(time.Hour * time.Duration(hoursInt)))}
		case days != "":
			daysInt, err := strconv.Atoi(days)
			if err != nil {
				return errors.New("days argument is not a valid number")
			}
			if daysInt <= 0 {
				return errors.New("days argument must be greater than 0")
			}
			query["insertionDate"] = bson.M{"$gt": time.Now().AddDate(0, 0, -daysInt)}
		case weeks != "":
			weeksInt, err := strconv.Atoi(weeks)
			if err != nil {
				return errors.New("weeks argument is not a valid number")
			}
			if weeksInt <= 0 {
				return errors.New("weeks argument must be greater than 0")
			}
			query["insertionDate"] = bson.M{"$gt": time.Now().AddDate(0, 0, -(7 * weeksInt))}
		}

		dbReferenceNrs := []models.ParsedCVReference{}
		err := dbConn.Find(&models.ParsedCVReference{}, &dbReferenceNrs, query)
		if err != nil {
			return err
		}

		res := make([]string, len(dbReferenceNrs))
		for idx, ref := range dbReferenceNrs {
			res[idx] = ref.ReferenceNumber
		}

		return c.JSON(res)
	},
}
