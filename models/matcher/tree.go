package matcher

import (
	"errors"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tree contains the cache for the tree so tree resolution can be fast
type Tree struct {
	branches     map[primitive.ObjectID]*Branch
	rootBranches []primitive.ObjectID
}

// GetBranch returns a branch with all the child branches in a tree structure
func (tc *Tree) GetBranch(dbConn db.Connection, branchID *primitive.ObjectID) (*Branch, error) {
	err := tc.build(dbConn)
	if err != nil {
		return nil, err
	}

	if branchID != nil {
		// TODO what do we do when a branch cannot be found?
		return tc.branches[*branchID], nil
	}

	parsedBranches := make([]*Branch, len(tc.rootBranches))
	for idx := range tc.rootBranches {
		// TODO what do we do when a branch cannot be found?
		parsedBranches[idx] = tc.branches[tc.rootBranches[idx]]
	}

	return &Branch{
		Titles:         nil,
		TitleKind:      Root,
		Branches:       tc.rootBranches,
		ParsedBranches: parsedBranches,
	}, nil
}

// GetIDsForBranch returns a spesific branches child branches their ids
func (tc *Tree) GetIDsForBranch(dbConn db.Connection, branchID primitive.ObjectID) ([]primitive.ObjectID, error) {
	err := tc.build(dbConn)
	if err != nil {
		return nil, err
	}

	resp := []primitive.ObjectID{}
	err = tc.findIDsForBranch(branchID, &resp)
	return resp, err
}

// findIDsForBranch adds a spesific branches child branches their ids to addTo parameter
func (tc *Tree) findIDsForBranch(branchID primitive.ObjectID, addTo *[]primitive.ObjectID) error {
	*addTo = append(*addTo, branchID)

	branch, ok := tc.branches[branchID]
	if !ok {
		return errors.New("branch not found")
	}

	for _, id := range branch.Branches {
		err := tc.findIDsForBranch(id, addTo)
		if err != nil {
			return err
		}
	}

	return nil
}

// build builds the tree cache
func (tc *Tree) build(dbConn db.Connection) error {
	branches := []Branch{}
	err := dbConn.Find(&Branch{}, &branches, bson.M{})
	if err != nil {
		return err
	}

	tc.branches = map[primitive.ObjectID]*Branch{}
	tc.rootBranches = []primitive.ObjectID{}

	// fill the branches map and set the HasParents property
	for idx := range branches {
		branch := branches[idx]

		for _, id := range branch.Branches {
			referencedBranch := tc.branches[id]
			if referencedBranch == nil {
				tc.branches[id] = &Branch{HasParents: true}
			} else {
				referencedBranch.HasParents = true
			}
		}

		if exsitingBranch, ok := tc.branches[branch.ID]; ok {
			branch.HasParents = exsitingBranch.HasParents
			*exsitingBranch = branch
		} else {
			tc.branches[branch.ID] = &branch
		}
	}

	// Link the branches to their parents
	// And find the root branches
	for _, branch := range tc.branches {
		branch.ParsedBranches = make([]*Branch, len(branch.Branches))
		for idx := 0; idx < len(branch.Branches); idx++ {
			// TODO what do we do when a branch cannot be found?
			referencedBranch := tc.branches[branch.Branches[idx]]
			branch.ParsedBranches[idx] = referencedBranch
		}
		if !branch.HasParents {
			tc.rootBranches = append(tc.rootBranches, branch.ID)
		}
	}

	return nil
}
