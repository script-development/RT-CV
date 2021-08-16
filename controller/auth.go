package controller

import "github.com/gofiber/fiber/v2"

func authRoutes(base fiber.Router) {
	auth := base.Group(`/auth`)
	auth.Get(`/salt`, func(c *fiber.Ctx) error {
		// TODO: Implement me
		return nil
	})
}
