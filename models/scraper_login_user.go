package models

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/nacl/box"
)

// ScraperLoginUsers defines all the users a scraper can use
type ScraperLoginUsers struct {
	db.M          `bson:",inline"`
	ScraperID     primitive.ObjectID `json:"scraperId" bson:"scraperId"`
	ScraperPubKey string             `json:"scraperPubKey" bson:"scraperPubKey"`
	Users         []ScraperLoginUser `json:"users"`
}

// ErrScraperNoPublicKey is used when you want to add a user to a scraper but the scraper has no public key yet used to encrypt the password
var ErrScraperNoPublicKey = errors.New("this scraper has no public key yet, add a public key to the scraper to store passwords")

// EncryptPassword encrypts a user password using the public key of a scraper
func (s *ScraperLoginUsers) EncryptPassword(password string) (string, error) {
	switch len(s.ScraperPubKey) {
	case 0:
		return "", ErrScraperNoPublicKey
	case 44:
		// This is what we expect, continue
	default:
		return "", errors.New("invalid scraper public key length")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(s.ScraperPubKey)
	if err != nil {
		return "", err
	}
	key := new([32]byte)
	copy(key[:], keyBytes)

	// We add 32 bits of random data at the start of the password to make it harder to decrypt
	dataToEncrypt := append(random.Generate32Bytes(), []byte(password)...)

	// We don't need to sign a message as the scrapers don't care about who sends the scraper users as long as it gets them
	// That's why we use SealAnonymous over Seal
	encrypted, err := box.SealAnonymous(nil, dataToEncrypt, key, crypto_rand.Reader)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// CollectionName should yield the collection name for the entry
func (*ScraperLoginUsers) CollectionName() string {
	return "scraperLoginUsers"
}

// Indexes implements db.Entry
func (*ScraperLoginUsers) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{Keys: bson.M{"scraperId": 1}},
	}
}

// ScraperLoginUser defines a user that can be used by a scraper to login into a scraped website
type ScraperLoginUser struct {
	Username          string `json:"username"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"encryptedPassword,omitempty" bson:"encryptedPassword"`
}
