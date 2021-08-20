package controller

import "github.com/gofiber/fiber/v2"

func authRoutes(base fiber.Router) {
	auth := base.Group(`/auth`)
	auth.Get(`/seed`, func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"seed": string(GetAuth(c).GetBaseSeed()),
		})
	})
}
