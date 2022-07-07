package models

import (
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Match contains information about a match
// We add omitempty to a lot of fields as it saves a lot of space in the database
type Match struct {
	RequestID   primitive.ObjectID      `json:"requestId" bson:"requestId"` // Maybe we should remove this one it adds minimal extra value
	ProfileID   primitive.ObjectID      `json:"profileId" bson:"profileId" description:"the profile this match was made with"`
	KeyID       primitive.ObjectID      `json:"keyId" bson:"keyId" description:"the key used to upload this CV, this will be the api key used by the scraper"`
	When        jsonHelpers.RFC3339Nano `json:"when"`
	ReferenceNr string                  `json:"referenceNr" bson:"referenceNr" description:"The reference number of the CV"`

	// Is this a debug match
	// This is currently only true if the match was made using the /tryMatcher dashboard page
	Debug bool `bson:",omitempty" json:"debug" description:"is this a debug match, this is currently only true if the match was made using the /tryMatcher dashboard page"`

	// The values below are non nil if a match was found
	// The result of the match is stored in the value of the field

	// The profile domain match that was found
	YearsSinceWork      *int `bson:",omitempty" json:"yearsSinceWork"`
	YearsSinceEducation *int `bson:",omitempty" json:"yearsSinceEducation"`
	// the education name of the profile that was matched
	Education *string `bson:",omitempty" json:"education"`
	// The profile desired profession match that was found
	DesiredProfession     *string              `bson:",omitempty" json:"desiredProfession"`
	ProfessionExperienced *string              `bson:",omitempty" json:"professionExperienced"`
	DriversLicense        bool                 `bson:",omitempty" json:"driversLicense"`
	ZipCode               *ProfileDutchZipcode `bson:",omitempty" json:"zipCode"`
}
