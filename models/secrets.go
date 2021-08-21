package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/script-development/RT-CV/db/dbInterfaces"
	"github.com/script-development/RT-CV/helpers/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Secret contains a secret value that can be stored in the database by a api user
// The secret value is encrypted with a key that is not stored on our side and is controlled by the api user
type Secret struct {
	dbInterfaces.M `bson:"inline"`
	KeyID          primitive.ObjectID `bson:"keyId"`
	Key            string
	Value          string
}

// CollectionName returns the collection name of a secret
func (*Secret) CollectionName() string {
	return "secrets"
}

// CreateSecret creates a secret
func CreateSecret(keyID primitive.ObjectID, key string, encryptionKey string, value []byte) (*Secret, error) {
	data, err := crypto.Encrypt(value, []byte(encryptionKey))
	if err != nil {
		return nil, err
	}
	if !json.Valid(value) {
		return nil, errors.New("expected json value")
	}
	return &Secret{
		M:     dbInterfaces.NewM(),
		KeyID: keyID,
		Key:   key,
		Value: base64.RawStdEncoding.EncodeToString(data),
	}, nil
}

// UnsafeMustCreateSecret Creates a secret and panics if an error is returned
func UnsafeMustCreateSecret(keyID primitive.ObjectID, key string, encryptionKey string, value []byte) *Secret {
	s, err := CreateSecret(keyID, key, encryptionKey, value)
	if err != nil {
		panic(err)
	}
	return s
}

// Decrypt decrypts the value of a secret
func (secret Secret) Decrypt(key string) (json.RawMessage, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(secret.Value)
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(bytes, []byte(key))
}

// GetSecretByKey gets a secret
func GetSecretByKey(conn dbInterfaces.Connection, keyID primitive.ObjectID, key string) (*Secret, error) {
	secret := &Secret{}
	err := conn.FindOne(secret, bson.M{"key": key, "keyId": keyID})
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// DeleteSecretByKey delete a secret
func DeleteSecretByKey(conn dbInterfaces.Connection, keyID primitive.ObjectID, key string) error {
	secret, err := GetSecretByKey(conn, keyID, key)
	if err != nil {
		return err
	}

	return conn.DeleteByID(secret)
}
