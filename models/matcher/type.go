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
	TitleKind TitleKind `bson:"titleKind" json:"titleKind"`

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

// GetBranch fetches a spesific branch from the database
func GetBranch(dbConn db.Connection, id primitive.ObjectID) (*Branch, error) {
	result := &Branch{}
	err := dbConn.FindOne(result, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetInTree returns all branches in a list in the tree or from a spesific point in the tree
func GetInTree(dbConn db.Connection, fromBranch *primitive.ObjectID) ([]Branch, error) {
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
	return branches, nil
}

// GetTree returns a db tree from the root or a spesific branch
func GetTree(dbConn db.Connection, fromBranch *primitive.ObjectID, deep int) (*Branch, error) {
	branches, err := GetInTree(dbConn, fromBranch)
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

	var idsToList func(branchIDs []primitive.ObjectID, nDeep int) []Branch
	idsToList = func(branchIDs []primitive.ObjectID, nDeep int) []Branch {
		if deep > 0 && nDeep == deep {
			return nil
		}
		result := []Branch{}
		for _, id := range branchIDs {
			branch, _ := branchesMap[id]
			branch.ParsedBranches = idsToList(branch.Branches, nDeep+1)
			result = append(result, branch)
		}
		return result
	}

	if fromBranch == nil {
		return &Branch{
			Titles:         nil,
			TitleKind:      Sector,
			ParsedBranches: idsToList(rootBrancheIDs, 1),
			Parents:        []primitive.ObjectID{},
		}, nil
	}

	return &idsToList(rootBrancheIDs, 0)[0], nil
}

// AddLeafProps are the leaf creation arguments for the AddLeaf method
type AddLeafProps struct {
	Titles    []string  `json:"titles"`
	TitleKind TitleKind `json:"titleKind"`
}

func (props AddLeafProps) validate() error {
	if len(props.Titles) == 0 {
		return errors.New("a title is required to add a leaf to a branch")
	}
	if props.TitleKind > 1 {
		return errors.New("titleKind invalid, must be 0 for a Job title and 1 for a Sector title")
	}
	return nil
}

// AddLeaf adds a new branch to
func (b *Branch) AddLeaf(dbConn db.Connection, props AddLeafProps, injectIntoSource bool) (*Branch, error) {
	err := props.validate()
	if err != nil {
		return nil, err
	}

	bIsRoot := b.ID.IsZero()

	newBranch := &Branch{
		M:              db.NewM(),
		Titles:         props.Titles,
		TitleKind:      props.TitleKind,
		Branches:       []primitive.ObjectID{},
		ParsedBranches: []Branch{},
		Parents:        []primitive.ObjectID{},
	}
	if !bIsRoot {
		// This is a sub branch of the tree, add the parents property
		newBranch.Parents = append(b.Parents, b.ID)
	}

	err = dbConn.Insert(newBranch)
	if err != nil {
		return nil, err
	}

	b.Branches = append(b.Branches, newBranch.ID)
	if injectIntoSource {
		b.ParsedBranches = append(b.ParsedBranches, *newBranch)
	}
	if !bIsRoot {
		err = dbConn.UpdateByID(b)
		if err != nil {
			return nil, err
		}
	}

	return newBranch, nil
}

// Update updates a spesific branches data
func (b *Branch) Update(dbConn db.Connection, props AddLeafProps) error {
	err := props.validate()
	if err != nil {
		return err
	}

	b.Titles = props.Titles
	b.TitleKind = props.TitleKind
	return dbConn.UpdateByID(b)
}
