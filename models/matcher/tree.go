package matcher

import (
	"errors"
	"sync"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*

Important:
As the Tree structure is a global public object it is not safe to use it concurrently unless the sync.Mutex is correctly locked
All public methods in this file should correctly lock and unlock the TreeCache it's mutex to make it concurrent safe

*/

// TreeCache contains the cache for the tree so tree resolution can be fast
type TreeCache struct {
	// a memory lock so we cannot get
	m sync.Mutex

	isBuild        bool
	branchBranches map[primitive.ObjectID][]primitive.ObjectID
	rootBranches   []primitive.ObjectID
}

// Tree contains the
var Tree = &TreeCache{}

// AddLeaf adds a leaf to the cache without resetting it :^)
func (tc *TreeCache) AddLeaf(branchID primitive.ObjectID, parentID *primitive.ObjectID) {
	tc.m.Lock()
	defer tc.m.Unlock()

	if !tc.isBuild {
		return
	}

	tc.rootBranches = append(tc.rootBranches, branchID)
	tc.branchBranches[branchID] = []primitive.ObjectID{}
	if parentID != nil && !parentID.IsZero() {
		tc.branchBranches[*parentID] = append(tc.branchBranches[*parentID], branchID)
	}
}

// Invalidate invalidates the tree cache and makes it so next request to a tree will rebuild it
func (tc *TreeCache) Invalidate() {
	tc.m.Lock()
	tc.branchBranches = nil
	tc.rootBranches = nil
	tc.isBuild = false
	tc.m.Unlock()
}

// GetBranch returns a branch with all the child branches in a tree structure
func (tc *TreeCache) GetBranch(dbConn db.Connection, branchID *primitive.ObjectID) (*Branch, error) {
	var allBranches []Branch
	if !tc.isBuild {
		// The cache is not yet build, lets rebuild it!
		var err error
		allBranches, err = tc.Rebuild(dbConn)
		if err != nil {
			return nil, err
		}
	}

	// Obtain the branches we need for the tree or a spesific branch of it
	if branchID == nil {
		if allBranches == nil {
			// If allBranches was not yet set by the cache rebuild, get the data from the database
			allBranches = []Branch{}
			err := dbConn.Find(&Branch{}, &allBranches, bson.M{})
			if err != nil {
				return nil, err
			}
		}
	} else {
		// Get the branch and its child branches ids we need to fetch
		ids := []primitive.ObjectID{}
		tc.m.Lock()
		err := tc.findIDsForBranch(*branchID, &ids)
		tc.m.Unlock()
		if err != nil {
			return nil, err
		}

		if allBranches == nil {
			allBranches = []Branch{}
			err := dbConn.Find(&Branch{}, &allBranches, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				return nil, err
			}
		}
	}

	branchesMap := branchesMap{}
	branchesMap.insert(allBranches...)
	branchesMap.link()

	if branchID != nil {
		return branchesMap.find(*branchID), nil
	}

	tc.m.Lock()
	defer tc.m.Unlock()

	parsedBranches := make([]*Branch, len(tc.rootBranches))
	for idx := range tc.rootBranches {
		parsedBranches[idx] = branchesMap.find(tc.rootBranches[idx])
	}

	return &Branch{
		Titles:         nil,
		TitleKind:      Root,
		Branches:       tc.rootBranches,
		ParsedBranches: parsedBranches,
	}, nil
}

// Rebuild builds the tree cache
func (tc *TreeCache) Rebuild(dbConn db.Connection) ([]Branch, error) {
	branches := []Branch{}
	err := dbConn.Find(&Branch{}, &branches, bson.M{})
	if err != nil {
		return nil, err
	}

	tc.m.Lock()
	defer tc.m.Unlock()

	tc.branchBranches = map[primitive.ObjectID][]primitive.ObjectID{}
	tc.rootBranches = []primitive.ObjectID{}

	referencesToBranch := map[primitive.ObjectID]uint16{}
	// fill the branches for branches
	for _, branch := range branches {
		tc.branchBranches[branch.ID] = branch.Branches
		for _, id := range branch.Branches {
			referencesToBranch[id]++
		}

		_, ok := referencesToBranch[branch.ID]
		if !ok {
			referencesToBranch[branch.ID] = 0
		}
	}

	// Find the root branches
	for id, count := range referencesToBranch {
		if count == 0 {
			tc.rootBranches = append(tc.rootBranches, id)
		}
	}

	tc.isBuild = true
	return branches, nil
}

// GetIDsForBranch returns a spesific branches child branches their ids
func (tc *TreeCache) GetIDsForBranch(dbConn db.Connection, branchID primitive.ObjectID) ([]primitive.ObjectID, error) {
	if !tc.isBuild {
		_, err := tc.Rebuild(dbConn)
		if err != nil {
			return nil, err
		}
	}

	resp := []primitive.ObjectID{}
	tc.m.Lock()
	err := tc.findIDsForBranch(branchID, &resp)
	tc.m.Unlock()
	return resp, err
}

// findIDsForBranch adds a spesific branches child branches their ids to ids
func (tc *TreeCache) findIDsForBranch(branchID primitive.ObjectID, addTo *[]primitive.ObjectID) error {
	*addTo = append(*addTo, branchID)

	branch, ok := tc.branchBranches[branchID]
	if !ok {
		return errors.New("branch not found")
	}

	for _, id := range branch {
		err := tc.findIDsForBranch(id, addTo)
		if err != nil {
			return err
		}
	}

	return nil
}

// branchesMap is a brancher map that can be used to quickly find a branch
type branchesMap [256][]*Branch

// insert inserts new branches into the branches map
func (bm *branchesMap) insert(branches ...Branch) {
	for idx, branch := range branches {
		id := branch.ID
		first := uint16(id[len(id)-1])
		bm[first] = append(bm[first], &branches[idx])
	}
}

// link sets the ParsedBranches for every branch in the map
func (bm *branchesMap) link() {
	for _, bucket := range bm {
		for _, branch := range bucket {
			branch.ParsedBranches = make([]*Branch, len(branch.Branches))
			for idx := 0; idx < len(branch.Branches); idx++ {
				// TODO what do we do when a branch cannot be found?
				branch.ParsedBranches[idx] = bm.find(branch.Branches[idx])
				bm.find(branch.Branches[idx])
			}
		}
	}
}

func (bm branchesMap) find(id primitive.ObjectID) *Branch {
	first := uint16(id[len(id)-1])
	potentialBranches := bm[first]

	for _, b := range potentialBranches {
		if b.ID == id {
			return b
		}
	}

	return nil
}
