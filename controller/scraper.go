package controller

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/match"
	"github.com/script-development/RT-CV/models"
)

func scraperRoutes(base fiber.Router) {
	scraper := base.Group(`/scraper`, requiresAuth(models.ApiKeyRoleScraper))
	scraper.Post(`/scanCV`, func(c *fiber.Ctx) error {
		body := models.Cv{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		profiles := GetProfiles(c)
		matchedProfiles := match.Match("werk.nl", *profiles, body) // TODO remove this hardcoded value
		if len(matchedProfiles) > 0 {
			for _, profile := range matchedProfiles {
				_, err := body.GetPDF(profile, "") // TODO add matchtext
				if err != nil {
					return fmt.Errorf("unable to generate PDF from CV, err: %s", err.Error())
				}
				fmt.Println(profile.Emails)
				// for _, email := range profile.Emails {
				// 	email.Email.Name
				// }
			}
		}

		return c.SendString("OK")
	})

	secret := scraper.Group(`/secret/:key`, validKey())
	secretKey := secret.Group(`/:secretKey`, validSecretKey())

	secret.Delete(``, func(c *fiber.Ctx) error {
		apiKey := GetKey(c)
		keyParam := c.Params("key")

		err := models.DeleteSecretByKey(GetDbConn(c), apiKey.ID, keyParam)
		if err != nil {
			return err
		}
		return c.JSON("ok")
	})
	secretKey.Post(``, func(c *fiber.Ctx) error {
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
	})
	secretKey.Put(``, func(c *fiber.Ctx) error {
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
	})
	secretKey.Get(``, func(c *fiber.Ctx) error {
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
	})
}

func validKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Params("key")) == 0 {
			return errors.New("key param cannot be empty")
		}
		return c.Next()
	}
}

func validSecretKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Params("secretKey")) < 16 {
			return errors.New("secretKey param must have a minimal length of 16 chars")
		}
		return c.Next()
	}
}
