package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Match contains information about a match
// We add omitempty to a lot of fields as it saves a lot of space in the database
type Match struct {
	db.M        `bson:",inline"`
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
	Domain              *string `bson:",omitempty" json:"domains"`
	YearsSinceWork      *int    `bson:",omitempty" json:"yearsSinceWork"`
	YearsSinceEducation *int    `bson:",omitempty" json:"yearsSinceEducation"`
	// the education name of the profile that was matched
	Education *string `bson:",omitempty" json:"education"`
	// The profile desired profession match that was found
	DesiredProfession     *string              `bson:",omitempty" json:"desiredProfession"`
	ProfessionExperienced *string              `bson:",omitempty" json:"professionExperienced"`
	DriversLicense        bool                 `bson:",omitempty" json:"driversLicense"`
	ZipCode               *ProfileDutchZipcode `bson:",omitempty" json:"zipCode"`
}

// CollectionName returns the collection name of the Profile
func (*Match) CollectionName() string {
	return "matches"
}

// Indexes implements db.Entry
func (*Match) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{Keys: bson.M{"profileId": 1}},
		{Keys: bson.M{"keyId": 1}},
		{Keys: bson.M{"when": 1}},
		{Keys: bson.M{"when": -1}},
		{Keys: bson.M{"referenceNr": 1}},
	}
}

// GetMatches returns all matches for a specific key
// If keyID is nil, all matches for all keys are returned
func GetMatches(dbConn db.Connection, keyID *primitive.ObjectID) ([]Match, error) {
	results := []Match{}
	if keyID != nil {
		err := dbConn.Find(&Match{}, &results, primitive.M{"keyId": keyID})
		return results, err
	}
	err := dbConn.Find(&Match{}, &results, primitive.M{})
	return results, err
}

// GetMatchesOnReferenceNr returns all matches that have been done on a ReferenceNr
func GetMatchesOnReferenceNr(dbConn db.Connection, referenceNr string, keyID *primitive.ObjectID) ([]Match, error) {
	query := bson.M{"referenceNr": referenceNr}
	if keyID != nil {
		query["keyId"] = keyID
	}

	results := []Match{}
	err := dbConn.Find(&Match{}, &results, query)
	return results, err
}

// GetMatchesSince returns all matches that have been done since a certain date+time
func GetMatchesSince(dbConn db.Connection, since time.Time, keyID *primitive.ObjectID) ([]Match, error) {
	query := bson.M{"when": bson.M{"$gt": since}}
	if keyID != nil {
		query["keyId"] = keyID
	}

	results := []Match{}
	err := dbConn.Find(&Match{}, &results, query)
	return results, err
}

// GetMatchSentence returns a
func (m *Match) GetMatchSentence() string {
	sentences := []string{}
	addReason := func(reason string) {
		sentences = append(sentences, reason)
	}

	if m.YearsSinceWork != nil {
		switch *m.YearsSinceWork {
		case 0:
			addReason("minder dan 1 jaar geleden sinds laatste werk ervaaring")
		case 1:
			addReason("1 jaar sinds laatste werk ervaaring")
		default:
			addReason(strconv.Itoa(*m.YearsSinceWork) + " jaren sinds laatste werk ervaaring")
		}
	}
	if m.YearsSinceEducation != nil {
		switch *m.YearsSinceEducation {
		case 0:
			addReason("minder dan 1 jaar sinds laatste opleiding")
		case 1:
			addReason("1 jaar sinds laatste opleiding")
		default:
			addReason(strconv.Itoa(*m.YearsSinceEducation) + " jaren sinds laatste opleiding")
		}
	}
	if m.Education != nil {
		addReason("opleiding " + *m.Education)
	}
	if m.DesiredProfession != nil {
		addReason("gewenste werkveld " + *m.DesiredProfession)
	}
	if m.ProfessionExperienced != nil {
		addReason("gewerkt als " + *m.ProfessionExperienced)
	}
	if m.DriversLicense {
		addReason("gewenste rijbewijs")
	}
	if m.ZipCode != nil {
		addReason(fmt.Sprintf("postcode in range %d - %d", m.ZipCode.From, m.ZipCode.To))
	}

	switch len(sentences) {
	case 0:
		return ""
	case 1:
		return sentences[0]
	default:
		return fmt.Sprintf("%s en %s", strings.Join(sentences[:len(sentences)-1], ", "), sentences[len(sentences)-1])
	}
}
