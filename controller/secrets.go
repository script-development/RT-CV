package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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

// TODO find a better name for this
func validSecretKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Params("secretKey")) < 16 {
			return errors.New("secretKey param must have a minimal length of 16 chars")
		}
		return c.Next()
	}
}

func routeCreateSecret(c *fiber.Ctx) error {
	apiKey := GetKey(c)
	keyParam, secretKeyParam := c.Params("key"), c.Params("secretKey")
	body := c.Body()
	if len(body) == 0 {
		return errors.New("body cannot be empty")
	}

	secret, err := models.CreateSecret(apiKey.ID, keyParam, secretKeyParam, body)
	if err != nil {
		return err
	}

	err = GetDbConn(c).Insert(secret)
	if err != nil {
		return err
	}

	secretValue, err := secret.Decrypt(secretKeyParam)
	if err != nil {
		return err
	}

	return c.JSON(secretValue)
}

func routeUpdateSecret(c *fiber.Ctx) error {
	apiKey := GetKey(c)
	keyParam, secretKeyParam := c.Params("key"), c.Params("secretKey")
	body := c.Body()
	if len(body) == 0 {
		return errors.New("body cannot be empty")
	}

	secret, err := models.GetSecretByKey(GetDbConn(c), apiKey.ID, keyParam)
	if err != nil {
		return err
	}
	// check if the key provided is equal to the previous key
	_, err = secret.Decrypt(secretKeyParam)
	if err != nil {
		return err
	}

	newSecret, err := models.CreateSecret(apiKey.ID, keyParam, secretKeyParam, body)
	if err != nil {
		return err
	}
	secret.Value = newSecret.Value

	// check if decryption still works
	secretValue, err := secret.Decrypt(secretKeyParam)
	if err != nil {
		return err
	}

	err = GetDbConn(c).UpdateByID(secret)
	if err != nil {
		return err
	}

	return c.JSON(secretValue)
}

func routeGetSecret(c *fiber.Ctx) error {
	apiKey := GetKey(c)
	keyParam, secretKeyParam := c.Params("key"), c.Params("secretKey")

	secret, err := models.GetSecretByKey(GetDbConn(c), apiKey.ID, keyParam)
	if err != nil {
		return err
	}

	value, err := secret.Decrypt(secretKeyParam)
	if err != nil {
		return err
	}

	return c.JSON(value)
}

func routeDeleteSecret(c *fiber.Ctx) error {
	apiKey := GetKey(c)
	keyParam := c.Params("key")

	err := models.DeleteSecretByKey(GetDbConn(c), apiKey.ID, keyParam)
	if err != nil {
		return err
	}
	return c.JSON("ok")
}
