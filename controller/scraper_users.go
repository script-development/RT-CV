package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func sendScraperLoginUsers(resp models.ScraperLoginUsers, userContext *ctx.Ctx, c *fiber.Ctx) error {
	if !userContext.Key.Roles.ContainsAll(models.APIKeyRoleScraper) {
		// Only the scrapers should be able to see the password of the login credentials
		// We copy the users slice here to prevent an issue with the testing db
		respUsers := make([]models.ScraperLoginUser, len(resp.Users))
		for idx, user := range resp.Users {
			user.Password = ""
			respUsers[idx] = user
		}
		resp.Users = respUsers
	}
	return c.JSON(resp)
}

var routeGetScraperUsers = routeBuilder.R{
	Description: "Get the login users of a specific scraper",
	Res:         models.ScraperLoginUsers{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)

		if !ctx.APIKeyFromParam.Roles.ContainsAll(models.APIKeyRoleScraper) {
			return errors.New("param key id must have the scraper role")
		}

		var resp models.ScraperLoginUsers
		err := ctx.DBConn.FindOne(&resp, bson.M{"scraperId": ctx.APIKeyFromParam.ID})
		if err == mongo.ErrNoDocuments {
			resp = models.ScraperLoginUsers{
				ScraperID: ctx.APIKeyFromParam.ID,
				Users:     []models.ScraperLoginUser{},
			}
		} else if err != nil {
			return err
		}

		return sendScraperLoginUsers(resp, ctx, c)
	},
}

// RouteDeleteScraperUserBody is the body of the routeDeleteScraperUser
type RouteDeleteScraperUserBody struct {
	Username string `json:"username"`
}

var routeDeleteScraperUser = routeBuilder.R{
	Description: "remove a scraper user from a scraper",
	Body:        RouteDeleteScraperUserBody{},
	Res:         models.ScraperLoginUsers{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)

		if !ctx.APIKeyFromParam.Roles.ContainsAll(models.APIKeyRoleScraper) {
			return errors.New("param key id must have the scraper role")
		}

		body := RouteDeleteScraperUserBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		scraperUsers := models.ScraperLoginUsers{}
		err = ctx.DBConn.FindOne(&scraperUsers, bson.M{"scraperId": ctx.APIKeyFromParam.ID})
		if err == mongo.ErrNoDocuments {
			return errors.New("username not found")
		}

		removedUser := false
		for idx := len(scraperUsers.Users) - 1; idx >= 0; idx-- {
			if scraperUsers.Users[idx].Username == body.Username {
				// Remove the user
				scraperUsers.Users = append(scraperUsers.Users[:idx], scraperUsers.Users[idx+1:]...)
				removedUser = true
				// Do not break here because if there are for some reason multiple users in the db with the same username we can remove them all
			}
		}
		if !removedUser {
			return errors.New("username not found")
		}

		err = ctx.DBConn.UpdateByID(&scraperUsers)
		if err != nil {
			return err
		}

		return sendScraperLoginUsers(scraperUsers, ctx, c)
	},
}

var routePatchScraperUser = routeBuilder.R{
	Description: "Update or insert a new scraper user",
	Body:        models.ScraperLoginUser{},
	Res:         models.ScraperLoginUsers{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)

		if !ctx.APIKeyFromParam.Roles.ContainsAll(models.APIKeyRoleScraper) {
			return errors.New("param key id must have the scraper role")
		}

		body := models.ScraperLoginUser{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}
		if body.Password == "" {
			return errors.New("password cannot be empty")
		}
		if body.Username == "" {
			return errors.New("username cannot be empty")
		}

		var alreadyExistingSet models.ScraperLoginUsers
		err = ctx.DBConn.FindOne(&alreadyExistingSet, bson.M{"scraperId": ctx.APIKeyFromParam.ID})
		if err == mongo.ErrNoDocuments {
			// Create a new entry
			resp := models.ScraperLoginUsers{
				M:         db.NewM(),
				ScraperID: ctx.APIKeyFromParam.ID,
				Users:     []models.ScraperLoginUser{body},
			}
			err = ctx.DBConn.Insert(&resp)
			if err != nil {
				return err
			}

			return sendScraperLoginUsers(resp, ctx, c)
		} else if err != nil {
			return err
		}

		// Update or insert the existing set
		// Check if the user already exists
		insert := true
		for idx, usr := range alreadyExistingSet.Users {
			if usr.Username == body.Username {
				insert = false
				alreadyExistingSet.Users[idx] = body
				break
			}
		}
		if insert {
			alreadyExistingSet.Users = append(alreadyExistingSet.Users, body)
		}

		err = ctx.DBConn.UpdateByID(&alreadyExistingSet)
		if err != nil {
			return err
		}

		return sendScraperLoginUsers(alreadyExistingSet, ctx, c)
	},
}
