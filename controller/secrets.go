package controller

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO find a better name for this
func validKeyMiddleware() routeBuilder.M {
	return routeBuilder.M{
		Fn: func(c *fiber.Ctx) error {
			if len(c.Params("key")) == 0 {
				return errors.New("key param cannot be empty")
			}
			return c.Next()
		},
	}
}

// RouteUpdateOrCreateSecret is the post data for the route below
type RouteUpdateOrCreateSecret struct {
	Value          json.RawMessage             `json:"value"`
	ValueStructure models.SecretValueStructure `json:"valueStructure"`
	Description    string                      `json:"description"`
	EncryptionKey  string                      `json:"encryptionKey"`
}

var routeUpdateOrCreateSecret = routeBuilder.R{
	Description: "Create or Update a secret for this specific api key and key combination.\n" +
		"note 1: we will never store the secret / encryption key on our side that's up to you.\n" +
		"note 2: the body must contain a valid json structure it doesn't matter what content",
	Res:  IMap{},
	Body: RouteUpdateOrCreateSecret{},
	Fn: func(c *fiber.Ctx) error {
		apiKey := ctx.GetAPIKeyFromParam(c)
		keyParam := c.Params("key")

		body := RouteUpdateOrCreateSecret{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		if len(body.EncryptionKey) < 16 {
			return errors.New("encryptionKey param must have a minimal length of 16 chars")
		}
		if !body.ValueStructure.Valid() {
			return errors.New("valueStructure does not contain a valid structure")
		}

		dbConn := ctx.GetDbConn(c)
		secret, err := models.GetSecretByKey(dbConn, apiKey.ID, keyParam)
		if err == mongo.ErrNoDocuments {
			// Secret does not yet exists, create it
			secret, err := models.CreateSecret(
				apiKey.ID,
				keyParam,
				body.EncryptionKey,
				body.Value,
				body.Description,
				body.ValueStructure,
			)
			if err != nil {
				return err
			}

			err = dbConn.Insert(secret)
			if err != nil {
				return err
			}

			secretValue, err := secret.Decrypt(body.EncryptionKey)
			if err != nil {
				return err
			}

			return c.JSON(secretValue)
		} else if err != nil {
			return err
		} else {
			// Secret exists

			// check if the key provided is equal to the previous key
			_, err = secret.Decrypt(body.EncryptionKey)
			if err != nil {
				return err
			}

			err := secret.UpdateValue(body.Value, body.EncryptionKey, body.ValueStructure)
			if err != nil {
				return err
			}
			secret.Description = body.Description

			// just making sure the decryption key still works
			secretValue, err := secret.Decrypt(body.EncryptionKey)
			if err != nil {
				return err
			}

			err = ctx.GetDbConn(c).UpdateByID(secret)
			if err != nil {
				return err
			}

			return c.JSON(secretValue)
		}
	},
}

var routeGetSecret = routeBuilder.R{
	Description: "Get a secret stored for this API Key and key pair",
	ResMap: map[string]interface{}{
		"free":  IMap{},
		"user":  models.SecretValueStructureUserT{},
		"users": []models.SecretValueStructureUserT{},
	},
	Fn: func(c *fiber.Ctx) error {
		dbConn := ctx.GetDbConn(c)
		apiKey := ctx.GetAPIKeyFromParam(c)
		keyParam, encryptionKeyParam := c.Params("key"), c.Params("encryptionKey")

		if len(encryptionKeyParam) < 16 {
			return errors.New("encryptionKey param must have a minimal length of 16 chars")
		}

		secret, err := models.GetSecretByKey(dbConn, apiKey.ID, keyParam)
		if err != nil {
			return err
		}

		value, err := secret.Decrypt(encryptionKeyParam)
		if err != nil {
			return err
		}

		return c.JSON(value)
	},
}

var routeGetSecrets = routeBuilder.R{
	Description: "Get all the stored secrets in the database for this specific API key " +
		"(this is without the secret value)",
	Res: []models.Secret{},
	Fn: func(c *fiber.Ctx) error {
		dbConn := ctx.GetDbConn(c)
		apiKey := ctx.GetAPIKeyFromParam(c)

		secrets, err := models.GetSecrets(dbConn, apiKey.ID)
		if err != nil {
			return err
		}

		return c.JSON(secrets)
	},
}

var routeGetAllSecretsFromAllKeys = routeBuilder.R{
	Description: "Get all the stored secrets in the database for all the API keys" +
		"(this is without the secret value)",
	Res: []models.Secret{},
	Fn: func(c *fiber.Ctx) error {
		dbConn := ctx.GetDbConn(c)

		secrets, err := models.GetSecretsFromAllKeys(dbConn)
		if err != nil {
			return err
		}
		return c.JSON(secrets)
	},
}

// RouteDeleteSecretOkRes is the response for the route below
type RouteDeleteSecretOkRes struct {
	Status string `json:"status"`
}

var routeDeleteSecret = routeBuilder.R{
	Description: "Delete a secret stored in the database",
	Res:         RouteDeleteSecretOkRes{},
	Fn: func(c *fiber.Ctx) error {
		apiKey := ctx.GetAPIKeyFromParam(c)
		keyParam := c.Params("key")

		err := models.DeleteSecretByKey(ctx.GetDbConn(c), apiKey.ID, keyParam)
		if err != nil {
			return err
		}
		return c.JSON(RouteDeleteSecretOkRes{"ok"})
	},
}
