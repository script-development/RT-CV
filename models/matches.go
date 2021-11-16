package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Match contains information about a match
// We add omitempty to a lot of fields as it saves a lot of space in the database
type Match struct {
	db.M      `bson:",inline"`
	RequestID primitive.ObjectID      `json:"requestId" bson:"requestId"` // Maybe we should remove this one it adds minimal extra value
	ProfileID primitive.ObjectID      `json:"profileId" bson:"profileId"`
	KeyID     primitive.ObjectID      `json:"keyId" bson:"keyId"`
	When      jsonHelpers.RFC3339Nano `json:"when"`

	// Is this a debug match
	// This is currently only true if the match was made using the /tryMatcher dashboard page
	Debug bool `bson:",omitempty" json:"debug"`

	// The values below are non nil if a match was found
	// The result of the match is stored in the value of the field

	// The profile domain match that was found
	Domain              *string `bson:",omitempty" json:"domains"`
	YearsSinceWork      *int    `bson:",omitempty" json:"yearsSinceWork"`
	YearsSinceEducation *int    `bson:",omitempty" json:"yearsSinceEducation"`
	// the education name of the profile that was matched
	Education *string `bson:",omitempty" json:"education"`
	// the course name of the profile that was matched
	// note that the profile does not really have a course we rather use the education name that matched a cv course
	Course *string `bson:",omitempty" json:"course"`
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
	if m.Course != nil {
		addReason("cursus " + *m.Course)
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
