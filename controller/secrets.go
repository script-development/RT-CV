package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

// TODO find a better name for this
func validKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Params("key")) == 0 {
			return errors.New("key param cannot be empty")
		}
		return c.Next()
	}
}

func validEncryptionKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Params("encryptionKey")) < 16 {
			return errors.New("encryptionKey param must have a minimal length of 16 chars")
		}
		return c.Next()
	}
}

func routeCreateSecret(c *fiber.Ctx) error {
	apiKey := ctx.GetAPIKeyFromParam(c)
	keyParam, encryptionKeyParam := c.Params("key"), c.Params("encryptionKey")
	body := c.Body()
	if len(body) == 0 {
		return errors.New("body cannot be empty")
	}

	secret, err := models.CreateSecret(apiKey.ID, keyParam, encryptionKeyParam, body)
	if err != nil {
		return err
	}

	err = ctx.GetDbConn(c).Insert(secret)
	if err != nil {
		return err
	}

	secretValue, err := secret.Decrypt(encryptionKeyParam)
	if err != nil {
		return err
	}

	return c.JSON(secretValue)
}

func routeUpdateSecret(c *fiber.Ctx) error {
	apiKey := ctx.GetAPIKeyFromParam(c)
	keyParam, encryptionKeyParam := c.Params("key"), c.Params("encryptionKey")
	body := c.Body()
	if len(body) == 0 {
		return errors.New("body cannot be empty")
	}

	secret, err := models.GetSecretByKey(ctx.GetDbConn(c), apiKey.ID, keyParam)
	if err != nil {
		return err
	}
	// check if the key provided is equal to the previous key
	_, err = secret.Decrypt(encryptionKeyParam)
	if err != nil {
		return err
	}

	newSecret, err := models.CreateSecret(apiKey.ID, keyParam, encryptionKeyParam, body)
	if err != nil {
		return err
	}
	secret.Value = newSecret.Value

	// check if decryption still works
	secretValue, err := secret.Decrypt(encryptionKeyParam)
	if err != nil {
		return err
	}

	err = ctx.GetDbConn(c).UpdateByID(secret)
	if err != nil {
		return err
	}

	return c.JSON(secretValue)
}

var routeGetSecret = routeBuilder.R{
	Description: "Get a secret stored for this API Key and key pair",
	Res:         IMap{},
	Fn: func(c *fiber.Ctx) error {
		dbConn := ctx.GetDbConn(c)
		apiKey := ctx.GetAPIKeyFromParam(c)
		keyParam, encryptionKeyParam := c.Params("key"), c.Params("encryptionKey")

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

func routeDeleteSecret(c *fiber.Ctx) error {
	apiKey := ctx.GetAPIKeyFromParam(c)
	keyParam := c.Params("key")

	err := models.DeleteSecretByKey(ctx.GetDbConn(c), apiKey.ID, keyParam)
	if err != nil {
		return err
	}
	return c.JSON(IMap{"status": "ok"})
}
