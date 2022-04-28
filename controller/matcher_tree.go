package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models/matcher"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeGetMatcherTree = routeBuilder.R{
	Description: "get the full matcher tree",
	Res:         matcher.Branch{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)

		tree, err := matcher.GetTree(ctx.DBConn, nil)
		if err != nil {
			return err
		}

		return c.JSON(tree)
	},
}

var routeGetPartOfMatcherTree = routeBuilder.R{
	Description: "get a spesific part of the matcher tree defined by it's id",
	Res:         matcher.Branch{},
	Fn: func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return err
		}

		ctx := ctx.Get(c)

		tree, err := matcher.GetTree(ctx.DBConn, &id)
		if err != nil {
			return err
		}

		return c.JSON(tree)
	},
}