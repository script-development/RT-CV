package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func routeGetKeys(c *fiber.Ctx) error {
	dbConn := ctx.GetDbConn(c)
	keys, err := models.GetAPIKeys(dbConn)
	if err != nil {
		return err
	}
	return c.JSON(keys)
}

type apiKeyModifyCreateData struct {
	Enabled *bool              `json:"enabled"`
	Domains []string           `json:"domains"`
	Key     *string            `json:"key"`
	Roles   *models.APIKeyRole `json:"roles"`
}

func routeCreateKey(c *fiber.Ctx) error {
	dbConn := ctx.GetDbConn(c)
	authenticator := ctx.GetAuth(c)

	body := apiKeyModifyCreateData{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	newAPIKey := &models.APIKey{
		M:       db.NewM(),
		Enabled: body.Enabled == nil || *body.Enabled,
	}

	if body.Domains == nil {
		return errors.New("domains should be set")
	} else if len(body.Domains) < 1 {
		return errors.New("there should at least be one domain")
	} else {
		// TODO check if each domain is valid
		newAPIKey.Domains = body.Domains
	}

	if body.Key == nil {
		return errors.New("key should be set")
	} else if len(*body.Key) < 16 {
		return errors.New("key must have a length of at least 16 chars")
	} else {
		newAPIKey.Key = *body.Key
	}

	if body.Roles == nil {
		return errors.New("roles should be set")
	} else if !body.Roles.Valid() {
		return errors.New("roles are invalid")
	} else {
		newAPIKey.Roles = *body.Roles
	}

	err = dbConn.Insert(newAPIKey)
	if err != nil {
		return err
	}

	authenticator.AddKey(*newAPIKey)

	return c.JSON(newAPIKey)
}

func routeDeleteKey(c *fiber.Ctx) error {
	apiKey := ctx.GetAPIKeyFromParam(c)
	if apiKey.System {
		return errors.New("you are not allowed to remove system keys")
	}

	dbConn := ctx.GetDbConn(c)
	err := dbConn.DeleteByID(apiKey)
	if err != nil {
		return err
	}
	return c.JSON(apiKey)
}

func routeGetKey(c *fiber.Ctx) error {
	apiKey := ctx.GetAPIKeyFromParam(c)
	return c.JSON(apiKey)
}

func routeUpdateKey(c *fiber.Ctx) error {
	dbConn := ctx.GetDbConn(c)
	apiKey := ctx.GetAPIKeyFromParam(c)
	authenticator := ctx.GetAuth(c)
	if apiKey.System {
		return errors.New("you are not allowed to remove system keys")
	}

	body := apiKeyModifyCreateData{}
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	if body.Enabled != nil {
		apiKey.Enabled = *body.Enabled
	}

	if body.Domains != nil {
		if len(body.Domains) < 1 {
			return errors.New("there should at least be one domain")
		}
		// TODO check if each domain is valid
		apiKey.Domains = body.Domains
	}

	keyChanged := false
	if body.Key != nil {
		if len(*body.Key) < 16 {
			return errors.New("key must have a length of at least 16 chars")
		}
		keyChanged = apiKey.Key != *body.Key
		apiKey.Key = *body.Key
	}

	if body.Roles != nil {
		if !body.Roles.Valid() {
			return errors.New("roles are invalid")
		}
		apiKey.Roles = *body.Roles
	}

	err = dbConn.UpdateByID(apiKey)
	if err != nil {
		return err
	}

	if keyChanged {
		authenticator.RefreshKey(*apiKey)
	}

	return c.JSON(apiKey)
}

func middlewareBindKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyParam := c.Params(`key`)
		keyID, err := primitive.ObjectIDFromHex(keyParam)
		if err != nil {
			return err
		}
		dbConn := ctx.GetDbConn(c)
		apiKey := models.APIKey{}
		err = dbConn.FindOne(&apiKey, bson.M{
			"_id": keyID,
		}, db.FindOptions{
			NoDefaultFilters: true,
		})
		if err != nil {
			return err
		}

		c.SetUserContext(
			ctx.SetAPIKeyFromParam(
				c.UserContext(),
				&apiKey,
			),
		)

		return c.Next()
	}
}
