package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

// SeedRes is the response data for auth seed
type SeedRes struct {
	Seed string `json:"seed"`
}

var routeAuthSeed = routeBuilder.R{
	Description: "Get the server seed, " +
		"This value is required in the authentication process and might change",
	Res: SeedRes{},
	Fn: func(c *fiber.Ctx) error {
		return c.JSON(SeedRes{
			Seed: string(ctx.GetAuth(c).GetBaseSeed()),
		})
	},
}

var routeGetKeyInfo = routeBuilder.R{
	Description: "Get information about the key you are using to authenticate with",
	Res:         models.APIKeyInfo{},
	Fn: func(c *fiber.Ctx) error {
		return c.JSON(ctx.GetKey(c).Info())
	},
}
