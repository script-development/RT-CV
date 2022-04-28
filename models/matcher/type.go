package matcher

import (
	"errors"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TitleKind defines the kind of the title
type TitleKind uint8

const (
	// Job name
	Job TitleKind = iota
	// Sector name
	Sector
)

// Branch contains a branch of the matcher tree
type Branch struct {
	db.M `bson:",inline"`

	// The titles of this branch
	// At least one title should be set
	Titles []string `json:"titles,omitempty"`

	// The kind of title
	// In the data we import there are 2 kinds of titles
	// One that contains a job name
	// And one that contains a Sector name
	TitleKind TitleKind `bson:"titleKind" json:"titleKind,omitempty"`

	// Branches contains sub branches ontop of this branch
	Branches []primitive.ObjectID `json:"-"`

	// ParsedBranches can be set when building a tree that is send to a user over the api in JSON format
	ParsedBranches []Branch `bson:"-" json:"branches,omitempty"`

	// Parents of this branch from the bottom of the tree up to this branch
	Parents []primitive.ObjectID `json:"parents,omitempty"`
}

// CollectionName implements db.Entry
func (*Branch) CollectionName() string {
	return "matcherBranches"
}

// Indexes implements db.Entry
func (*Branch) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{Keys: bson.M{"parents": 1}},
		{Keys: bson.M{"titles": 1}},
	}
}

// GetTree returns a db tree from the root or a spesific branch
func GetTree(dbConn db.Connection, fromBranch *primitive.ObjectID) (*Branch, error) {
	query := bson.M{}
	if fromBranch != nil {
		query["$or"] = []bson.M{
			{"parents": bson.M{"$exists": true, "$in": []any{*fromBranch}}},
			{"_id": fromBranch},
		}
	}

	branches := []Branch{}
	err := dbConn.Find(&Branch{}, &branches, query)
	if err != nil {
		return nil, err
	}

	rootBrancheIDs := []primitive.ObjectID{}
	branchesMap := map[primitive.ObjectID]Branch{}
	for _, v := range branches {
		branchesMap[v.ID] = v
		if fromBranch == nil {
			if len(v.Parents) == 0 {
				rootBrancheIDs = append(rootBrancheIDs, v.ID)
			}
		} else {
			if v.ID == *fromBranch {
				rootBrancheIDs = append(rootBrancheIDs, v.ID)
			}
		}
	}

	if len(rootBrancheIDs) == 0 {
		if fromBranch == nil {
			return nil, errors.New("no branches found")
		}
		return nil, errors.New("no branches for spesific branch found")
	}

	var idsToList func(branchIDs []primitive.ObjectID, path []primitive.ObjectID) []Branch
	idsToList = func(branchIDs []primitive.ObjectID, path []primitive.ObjectID) []Branch {
		result := []Branch{}
		for _, id := range branchIDs {
			branch, _ := branchesMap[id]
			branch.ParsedBranches = idsToList(branch.Branches, append(path, id))
			result = append(result, branch)
		}
		return result
	}

	parsedBranches := idsToList(rootBrancheIDs, []primitive.ObjectID{})
	if fromBranch == nil {
		return &Branch{
			Titles:         nil,
			TitleKind:      Sector,
			ParsedBranches: parsedBranches,
			Parents:        []primitive.ObjectID{},
		}, nil
	}

	return &parsedBranches[0], nil
}