package controller

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models/matcher"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var routeGetMatcherTree = routeBuilder.R{
	Description: "get the matcher tree or a spesific part of the matcher tree defined by it's id",
	Res:         matcher.Branch{},
	Fn: func(c *fiber.Ctx) error {
		var branchID *primitive.ObjectID
		idParam := c.Params("id")
		if idParam != "" {
			id, err := primitive.ObjectIDFromHex(idParam)
			if err != nil {
				return err
			}
			branchID = &id
		}

		resp := c.Response()
		resp.Header.SetContentType(fiber.MIMEApplicationJSON)

		fileInfo, file, err := matcher.ObtainCachedTree(branchID)
		if err == nil {
			// Yay we have a cache hit, lets return the cached response
			resp.SetBodyStream(file, int(fileInfo.Size()))
			return nil
		}

		ctx := ctx.Get(c)

		if os.IsNotExist(err) {
			ctx.Logger.Info("matcher tree is not cached (yet) for this id")
		} else {
			ctx.Logger.WithError(err).Warn("obtaining cached matcher tree failed, falling back to database")
		}

		// deep := c.Context().QueryArgs().GetUintOrZero("deep")
		tree, err := (&matcher.Tree{}).GetBranch(ctx.DBConn, branchID)
		if err != nil {
			return err
		}

		jsonTree, err := json.Marshal(tree)
		if err != nil {
			return err
		}
		err = matcher.CacheTree(branchID, jsonTree)
		if err != nil {
			return err
		}

		resp.SetBodyRaw(jsonTree)
		return nil
	},
}

// SearchMatcherLeafsRequest is the request for the search matcher leafs route
type SearchMatcherLeafsRequest struct {
	Search string `json:"search"`
}

var routeSearchMatcherLeaf = routeBuilder.R{
	Description: "search for a matcher leaf",
	Res:         []matcher.SearchResult{},
	Body:        SearchMatcherLeafsRequest{},
	Fn: func(c *fiber.Ctx) error {
		body := SearchMatcherLeafsRequest{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		ctx := ctx.Get(c)

		matches, err := (&matcher.Tree{}).Search(ctx.Logger, ctx.DBConn, body.Search)
		if err != nil {
			return err
		}

		return c.JSON(matches)
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

		defer matcher.NukeCache()

		_, err = tree.AddLeaf(ctx.DBConn, body, true /* deep != 1*/)
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

		defer matcher.NukeCache()

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
	Res:         RouteDeleteMatcherBranchResult{},
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

		defer matcher.NukeCache()

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
