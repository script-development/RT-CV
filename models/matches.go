package models

import (
	"fmt"
	"strings"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Match contains information about a match
// We add omitempty to a lot of fields as it saves a lot of space in the database
type Match struct {
	db.M      `bson:",inline"`
	RequestID primitive.ObjectID      `json:"requestId" bson:"requestId"`
	ProfileID primitive.ObjectID      `json:"profileId" bson:"profileId"`
	KeyID     primitive.ObjectID      `json:"keyId" bson:"keyId"`
	When      jsonHelpers.RFC3339Nano `json:"when"`

	Debug bool `bson:",omitempty" json:"debug"`

	// The profile domain match that was found
	Domain                *string              `bson:",omitempty" json:"domains"`
	YearsSinceWork        bool                 `bson:",omitempty" json:"yearsSinceWork"`
	YearsSinceEducation   bool                 `bson:",omitempty" json:"yearsSinceEducation"`
	EducationOrCourse     bool                 `bson:",omitempty" json:"educationOrCourse"`
	DesiredProfession     bool                 `bson:",omitempty" json:"desiredProfession"`
	ProfessionExperienced bool                 `bson:",omitempty" json:"professionExperienced"`
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
	if m.Domain != nil {
		sentences = append(sentences, "domain naam "+*m.Domain)
	}
	if m.YearsSinceWork {
		sentences = append(sentences, "jaren sinds werk")
	}
	if m.YearsSinceEducation {
		sentences = append(sentences, "jaren sinds laatste opleiding")
	}
	if m.EducationOrCourse {
		sentences = append(sentences, "opleiding of cursus")
	}
	if m.DesiredProfession {
		sentences = append(sentences, "gewenste werkveld")
	}
	if m.ProfessionExperienced {
		sentences = append(sentences, "gewenst beroep")
	}
	if m.DriversLicense {
		sentences = append(sentences, "rijbewijs")
	}
	if m.ZipCode != nil {
		sentences = append(sentences, fmt.Sprintf("postcode in range %d - %d", m.ZipCode.From, m.ZipCode.To))
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
