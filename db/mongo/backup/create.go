package backup

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/dbHelpers"
	"github.com/script-development/RT-CV/db/mongo"
	"github.com/script-development/RT-CV/helpers/crypto"
	"github.com/script-development/RT-CV/helpers/numbers"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson"
)

/*

Data layout created by backup method

crypto.Encrypt(
	gzip({collections data})
	master_key,
)

The {collections data} is as follows, for every collection this is added:
[]byte{
	...name_length as uint16 (2 bytes),
	...name,
	...{collection data},
}

The {collection data} is as follows, every document is added like this:
[]byte{
	...raw_bson_data_length as uint64 (8 bytes),
	...raw_bson_data,
	{isLast bool (uint8 either 1 = true or 0 = false)},
}
*/

// CreateBackupFile creates a backup file from the database contents
//
// YOU NEED TO CLOSE THE RETURNED FILE
func CreateBackupFile(genericConn db.Connection, masterKey string) (*os.File, error) {
	backupFile, err := createBackupWriter(masterKey, func(w io.Writer) error {
		conn, ok := genericConn.(*mongo.Connection)
		if !ok {
			return errors.New("DB Connection is not a Mongo DB connection")
		}
		db := conn.GetDB()

		names, err := db.ListCollectionNames(dbHelpers.Ctx(), bson.M{})
		if err != nil {
			return err
		}

		for _, name := range names {
			ctx := dbHelpers.Ctx()
			cursor, err := db.Collection(name).Find(ctx, bson.M{})
			if err != nil {
				return err
			}
			first := true
			for cursor.Next(ctx) {
				document := make(bson.Raw, len(cursor.Current))
				copy(document, cursor.Current)
				if first {
					// Only write the name of the collection once we are sure
					collectionNameData := append(numbers.UintToBytes(uint64(len(name)), 16), []byte(name)...)
					w.Write(collectionNameData)

					first = false
				} else {
					// Write the is last document byte
					// In this case there is a next document so we write a false / 0
					w.Write([]byte{0})
				}

				w.Write(numbers.UintToBytes(uint64(len(document)), 64))
				w.Write(document)
			}
			if !first {
				// Write the last document byte
				// Only write the last document identifier if there where actually documents in this collection
				w.Write([]byte{1})
			}
			cursor.Close(ctx)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Info("validating generated backup file..")

	// Validate the generated data
	err = readbackup(backupFile, masterKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate backup data: %s", err)
	}

	_, err = backupFile.Seek(0, 0)
	if err != nil {
		backupFile.Close()
		os.Remove(backupFile.Name())

		return nil, err
	}

	err = models.SetLastBackupToNow(genericConn)
	if err != nil {
		backupFile.Close()
		os.Remove(backupFile.Name())

		return nil, err
	}

	return backupFile, nil
}

// createBackupWriter creates all the writers and closes them correctly in order in case of a error
func createBackupWriter(masterKey string, createBackupMethod func(io.Writer) error) (*os.File, error) {
	backupFile, err := os.Create("./backup.gz.aes")
	if err != nil {
		return nil, err
	}

	closeAndDeleteFile := func() {
		backupFile.Close()
		os.Remove(backupFile.Name())
	}

	encryptionWriter, err := crypto.NewEncryptWriter([]byte(masterKey), backupFile)
	if err != nil {
		closeAndDeleteFile()
		return nil, err
	}

	zw, err := gzip.NewWriterLevel(encryptionWriter, 5)
	if err != nil {
		encryptionWriter.Close()
		closeAndDeleteFile()
		return nil, err
	}

	err = createBackupMethod(zw)

	// Close writers
	zw.Close()
	encryptionWriter.Close()
	if err != nil {
		closeAndDeleteFile()
		return nil, err
	}

	err = backupFile.Sync()
	if err != nil {
		closeAndDeleteFile()
		return nil, err
	}

	_, err = backupFile.Seek(0, 0)
	if err != nil {
		closeAndDeleteFile()
		return nil, err
	}

	return backupFile, nil
}
