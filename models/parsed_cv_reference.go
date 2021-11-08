package models

import (
	"time"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReferenceNrIsParsed retruns if the reference number is parsed
func ReferenceNrIsParsed(conn db.Connection, keyID primitive.ObjectID, referenceNr string) bool {
	err := conn.FindOne(&ParsedCVReference{}, bson.M{"keyId": keyID, "referenceNumber": referenceNr})
	return err != nil
}

// InsertParsedCVReference inserts a new ParsedCVReference into the database
func InsertParsedCVReference(conn db.Connection, keyID primitive.ObjectID, referenceNr string) error {
	newEntry := &ParsedCVReference{
		M:               db.NewM(),
		ReferenceNumber: referenceNr,
		InsertionDate:   jsonHelpers.RFC3339Nano(time.Now()),
		KeyID:           keyID,
	}
	err := conn.Insert(newEntry)
	if err != nil {
		return err
	}

	return nil
}

// ParsedCVReference is a entry in the database that is used to detect duplicates in uploaded CVs
type ParsedCVReference struct {
	db.M            `bson:",inline"`
	ReferenceNumber string                  `bson:"referenceNumber" json:"referenceNumber"`
	InsertionDate   jsonHelpers.RFC3339Nano `bson:"insertionDate" json:"insertionDate"`
	KeyID           primitive.ObjectID      `bson:"keyId" json:"keyId"`
}

// CollectionName returns the collection name of the ParsedCVReference
func (*ParsedCVReference) CollectionName() string {
	return "parsedCvReferences"
}
