package controller

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/helpers/validation"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeGetScraperKeys = routeBuilder.R{
	Description: "Get all scraper keys from the database",
	Res:         []models.APIKey{},
	Fn: func(c *fiber.Ctx) error {
		keys, err := models.GetScraperAPIKeys(ctx.Get(c).DBConn)
		if err != nil {
			return err
		}
		return c.JSON(keys)
	},
}

var routeGetKeys = routeBuilder.R{
	Description: "get all api keys from the database",
	Res:         []models.APIKey{},
	Fn: func(c *fiber.Ctx) error {
		keys, err := models.GetAPIKeys(ctx.Get(c).DBConn)
		if err != nil {
			return err
		}
		return c.JSON(keys)
	},
}

type apiKeyModifyCreateData struct {
	Enabled *bool              `json:"enabled"`
	Name    *string            `json:"name"`
	Domains []string           `json:"domains"`
	Key     *string            `json:"key"`
	Roles   *models.APIKeyRole `json:"roles"`
}

var routeCreateKey = routeBuilder.R{
	Description: "create a new api key",
	Body:        apiKeyModifyCreateData{},
	Res:         models.APIKey{},
	Fn: func(c *fiber.Ctx) error {
		body := apiKeyModifyCreateData{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		newAPIKey := &models.APIKey{
			M:       db.NewM(),
			Enabled: body.Enabled == nil || *body.Enabled,
		}

		if body.Name == nil {
			return errors.New("name is required")
		} else if len(*body.Name) == 0 {
			return errors.New("name cannot be empty")
		}
		newAPIKey.Name = *body.Name

		if body.Domains == nil {
			return errors.New("domains should be set")
		} else if len(body.Domains) < 1 {
			return errors.New("there should at least be one domain")
		} else {
			err := validation.ValidDomainListAndFormat(&body.Domains, true)
			if err != nil {
				return err
			}
			for idx, domain := range body.Domains {
				body.Domains[idx] = strings.ToLower(domain)
			}
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

		err = ctx.Get(c).DBConn.Insert(newAPIKey)
		if err != nil {
			return err
		}

		return c.JSON(newAPIKey)
	},
}

var routeDeleteKey = routeBuilder.R{
	Description: "delete an api key",
	Res:         models.APIKey{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)
		if ctx.APIKeyFromParam.System {
			return errors.New("you are not allowed to remove system keys")
		}

		err := ctx.DBConn.DeleteByID(&models.APIKey{}, ctx.APIKeyFromParam.ID)
		if err != nil {
			return err
		}

		ctx.Auth.RemoveKeyCache(ctx.APIKeyFromParam.ID.Hex())

		return c.JSON(ctx.APIKeyFromParam)
	},
}

var routeGetKey = routeBuilder.R{
	Description: "get an api key from the database based on it's ID",
	Res:         models.APIKey{},
	Fn: func(c *fiber.Ctx) error {
		return c.JSON(ctx.Get(c).APIKeyFromParam)
	},
}

var routeUpdateKey = routeBuilder.R{
	Description: "Update an api key",
	Body:        apiKeyModifyCreateData{},
	Res:         models.APIKey{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)
		apiKey := ctx.APIKeyFromParam
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

		if body.Name != nil {
			apiKey.Name = *body.Name
		}

		if body.Domains != nil {
			if len(body.Domains) < 1 {
				return errors.New("there should at least be one domain")
			}
			err := validation.ValidDomainListAndFormat(&body.Domains, true)
			if err != nil {
				return err
			}
			for idx, domain := range body.Domains {
				body.Domains[idx] = strings.ToLower(domain)
			}
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

		err = ctx.DBConn.UpdateByID(apiKey)
		if err != nil {
			return err
		}

		if keyChanged {
			ctx.Auth.RemoveKeyCache(apiKey.ID.Hex())
		}

		return c.JSON(apiKey)
	},
}

// middlewareBindMyKey sets the APIKeyFromParam to the api key used to authenticate
func middlewareBindMyKey() routeBuilder.M {
	return routeBuilder.M{
		Fn: func(c *fiber.Ctx) error {
			ctx := ctx.Get(c)
			ctx.APIKeyFromParam = ctx.Key
			return c.Next()
		},
	}
}

func middlewareBindKey(urlParamName string) routeBuilder.M {
	return routeBuilder.M{
		Fn: func(c *fiber.Ctx) error {
			keyParam := c.Params(urlParamName)
			keyID, err := primitive.ObjectIDFromHex(keyParam)
			if err != nil {
				return err
			}
			ctx := ctx.Get(c)
			apiKey := models.APIKey{}
			query := bson.M{"_id": keyID}
			args := db.FindOptions{NoDefaultFilters: true}
			err = ctx.DBConn.FindOne(&apiKey, query, args)
			if err != nil {
				return err
			}

			ctx.APIKeyFromParam = &apiKey

			return c.Next()
		},
	}
}
