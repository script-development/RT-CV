package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models/matcher"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeGetMatcherTree = routeBuilder.R{
	Description: "get the full matcher tree",
	Res:         matcher.Branch{},
	Fn: func(c *fiber.Ctx) error {
		ctx := ctx.Get(c)

		tree, err := matcher.GetTree(ctx.DBConn, nil, c.Context().QueryArgs().GetUintOrZero("deep"))
		if err != nil {
			return err
		}

		return c.JSON(tree)
	},
}

var routeAddMatcherLeaf = routeBuilder.R{
	Description: "add a leaf to a spesific part of the tree",
	Res:         matcher.Branch{},
	Body:        matcher.AddLeafProps{},
	Fn: func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return err
		}

		body := matcher.AddLeafProps{}
		err = c.BodyParser(&body)
		if err != nil {
			return err
		}

		ctx := ctx.Get(c)

		tree, err := matcher.GetTree(ctx.DBConn, &id, c.Context().QueryArgs().GetUintOrZero("deep"))
		if err != nil {
			return err
		}

		_, err = tree.AddLeaf(ctx.DBConn, body)
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

		tree, err := matcher.GetTree(ctx.DBConn, &id, c.Context().QueryArgs().GetUintOrZero("deep"))
		if err != nil {
			return err
		}

		return c.JSON(tree)
	},
}

var routePutMatcherBranch = routeBuilder.R{
	Description: `update a spesific branch`,
	Body:        matcher.AddLeafProps{},
	Res:         matcher.Branch{},
	Fn: func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return err
		}

		body := matcher.AddLeafProps{}
		err = c.BodyParser(&body)
		if err != nil {
			return err
		}

		ctx := ctx.Get(c)

		tree, err := matcher.GetTree(ctx.DBConn, &id, c.Context().QueryArgs().GetUintOrZero("deep"))
		if err != nil {
			return err
		}

		err = tree.Update(ctx.DBConn, body)
		if err != nil {
			return err
		}

		return c.JSON(tree)
	},
}

var routeDeleteMatcherBranch = routeBuilder.R{
	Description: `remove a spesific branch and it's children`,
	Res:         RouteDeleteSecretOkRes{},
	Fn: func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return err
		}

		ctx := ctx.Get(c)

		branches, err := matcher.GetInTree(ctx.DBConn, &id)
		if err != nil {
			return err
		}
		if len(branches) == 0 {
			return errors.New("branch not found")
		}
		var parentID primitive.ObjectID
		found := false
		for _, v := range branches {
			if v.ID == id {
				parentsLen := len(v.Parents)
				if parentsLen != 0 {
					parentID = v.Parents[parentsLen-1]
				}
				found = true
				break
			}
		}
		if !found {
			return errors.New("branch not found")
		}

		if !parentID.IsZero() {
			parent, err := matcher.GetBranch(ctx.DBConn, parentID)
			if err != nil {
				return errors.New("parent of branch not found")
			}
			// Remove in reverse the ID from this branch as it no longer exsists
			for idx := len(parent.Branches) - 1; idx >= 0; idx-- {
				if parent.Branches[idx] == id {
					parent.Branches = append(parent.Branches[:idx], parent.Branches[idx+1:]...)
				}
			}
			err = ctx.DBConn.UpdateByID(parent)
			if err != nil {
				return err
			}
		}

		branchesToRemove := []db.Entry{}
		for idx := range branches {
			branchesToRemove = append(branchesToRemove, &branches[idx])
		}

		err = ctx.DBConn.DeleteByID(branchesToRemove...)
		if err != nil {
			return err
		}

		return c.JSON(RouteDeleteSecretOkRes{"ok"})
	},
}
