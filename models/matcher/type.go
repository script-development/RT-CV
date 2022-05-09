package matcher

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/mjarkk/jsonschema"
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TitleKind defines the kind of the title
type TitleKind uint8

// JSONSchemaDescribe implements jsonschema.Describe
func (TitleKind) JSONSchemaDescribe() jsonschema.Property {
	describe := `0 for a job title
1 for a sector title
2 for a education title
3 for the root of the tree`

	return jsonschema.Property{
		Title:       "The title kind",
		Description: describe,
		Type:        jsonschema.PropertyTypeInteger,
		Enum: []json.RawMessage{
			Job.toJSON(),
			Sector.toJSON(),
			Education.toJSON(),
			Root.toJSON(),
		},
	}
}

// Valid returns an error if the tree is valid
func (k TitleKind) Valid() error {
	// Note that Root is special and should not be used in the tree
	if k >= Job && k <= Education {
		return nil
	}
	return errors.New("titleKind not valid")
}

func (k TitleKind) toJSON() json.RawMessage {
	return []byte(strconv.FormatUint(uint64(k), 10))
}

const (
	// Job name
	Job TitleKind = iota
	// Sector name
	Sector
	// Education name
	Education
	// Root of the tree
	Root
)

// Branch contains a branch of the matcher tree
type Branch struct {
	db.M `bson:",inline"`

	// The titles of this branch
	// At least one title should be set
	Titles []string `bson:"title" json:"titles,omitempty"`

	// The kind of title
	// In the data we import there are 2 kinds of titles
	// One that contains a job name
	// And one that contains a Sector name
	TitleKind TitleKind `bson:"titleKind" json:"titleKind"`

	// Branches contains sub branches ontop of this branch
	Branches []primitive.ObjectID `json:"branchesIds,omitempty"`

	// ParsedBranches can be set when building a tree that is send to a user over the api in JSON format
	ParsedBranches []*Branch `bson:"-" json:"branches,omitempty"`

	// Used by the tree to find the root branches
	HasParents bool `json:"-" bson:"-"`
}

// CollectionName implements db.Entry
func (*Branch) CollectionName() string {
	return "matcherBranches"
}

// Indexes implements db.Entry
func (*Branch) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
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

// FindParents find the branchs that have the id arg as parent
func FindParents(dbConn db.Connection, id primitive.ObjectID) ([]Branch, error) {
	results := []Branch{}
	err := dbConn.Find(&Branch{}, &results, bson.M{"branches": id})
	return results, err
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
	return props.TitleKind.Valid()
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
		ParsedBranches: []*Branch{},
	}

	err = dbConn.Insert(newBranch)
	if err != nil {
		return nil, err
	}

	b.Branches = append(b.Branches, newBranch.ID)
	if injectIntoSource {
		b.ParsedBranches = append(b.ParsedBranches, newBranch)
	}
	if !bIsRoot {
		err = dbConn.UpdateByID(b)
		if err != nil {
			return nil, err
		}
	}

	err = NukeCache()
	if err != nil {
		return nil, err
	}

	return newBranch, nil
}

// Update updates a spesific branches data
func (b *Branch) Update(dbConn db.Connection, props AddLeafProps) error {
	err := props.validate()
	if err != nil {
		return err
	}

	err = NukeCache()
	if err != nil {
		return err
	}

	b.Titles = props.Titles
	b.TitleKind = props.TitleKind
	return dbConn.UpdateByID(b)
}
