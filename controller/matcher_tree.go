package controller

import (
	"errors"

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

		// deep := c.Context().QueryArgs().GetUintOrZero("deep")
		tree, err := (&matcher.Tree{}).GetBranch(ctx.DBConn, nil)
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
		body := matcher.AddLeafProps{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		ctx := ctx.Get(c)

		// deep := c.Context().QueryArgs().GetUintOrZero("deep")

		idParamStr := c.Params("id")
		var idParam *primitive.ObjectID
		if idParamStr != "" {
			id, err := primitive.ObjectIDFromHex(idParamStr)
			if err != nil {
				return err
			}
			idParam = &id
		}
		tree, err := (&matcher.Tree{}).GetBranch(ctx.DBConn, idParam)
		if err != nil {
			return err
		}

		_, err = tree.AddLeaf(ctx.DBConn, body, true /* deep != 1*/)
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

		// deep := c.Context().QueryArgs().GetUintOrZero("deep")
		tree, err := (&matcher.Tree{}).GetBranch(ctx.DBConn, &id)
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

		// deep := c.Context().QueryArgs().GetUintOrZero("deep")
		tree, err := (&matcher.Tree{}).GetBranch(ctx.DBConn, &id)
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

// RouteDeleteMatcherBranchResult gives some information about the removal process
type RouteDeleteMatcherBranchResult struct {
	UpdatedParents  int `json:"updatedParent"`
	DeletedBranches int `json:"deletedBranches"`
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

		// Firstly lets remove the parents as if this fails the database isn't broken and if the remaining code fails at least the data won't show up annymore
		parents, err := matcher.FindParents(ctx.DBConn, id)
		if err != nil {
			return err
		}
		for _, parent := range parents {
			for idx, branchID := range parent.Branches {
				if branchID == id {
					parent.Branches = append(parent.Branches[:idx], parent.Branches[idx+1:]...)
					break
				}
			}
			err = ctx.DBConn.UpdateByID(&parent)
			if err != nil {
				return err
			}
		}

		// delete the branch and all it's child branches
		branchIDs, err := (&matcher.Tree{}).GetIDsForBranch(ctx.DBConn, id)
		if err != nil {
			return err
		}
		if len(branchIDs) == 0 {
			return errors.New("branch not found")
		}
		err = ctx.DBConn.DeleteByID(&matcher.Branch{}, branchIDs...)
		if err != nil {
			return err
		}

		return c.JSON(RouteDeleteMatcherBranchResult{
			UpdatedParents:  len(parents),
			DeletedBranches: len(branchIDs),
		})
	},
}
