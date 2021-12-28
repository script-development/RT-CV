package mongo

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/dbHelpers"
	"github.com/script-development/RT-CV/helpers/crypto"
	"github.com/script-development/RT-CV/helpers/numbers"
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

func readbackup(fileReader io.Reader, masterKey string, restoreMethod func(collection string, doc bson.Raw, err error)) error {
	cryptoReader, err := crypto.NewEncryptReader([]byte(masterKey), fileReader)
	if err != nil {
		return err
	}

	zr, err := gzip.NewReader(cryptoReader)
	if err != nil {
		return err
	}

	mustRead := func(nr uint64, what string) ([]byte, error) {
		buf := make([]byte, nr)
		_, err := io.ReadFull(zr, buf)
		if err != nil {
			return nil, fmt.Errorf("failed to read bytes for %s: %s", what, err.Error())
		}
		return buf, nil
	}
	mustReadUint := func(size uint8, what string) (uint64, error) {
		buf, err := mustRead(uint64(size)/8, what)
		if err != nil {
			return 0, err
		}
		nr, err := numbers.BytesToUint(buf)
		if err != nil {
			return 0, fmt.Errorf("failed to convert bytes to uint for %s: %s", what, err.Error())
		}
		return nr, nil
	}

	for {
		buf := []byte{0, 0}
		bytesRead, err := zr.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Reached last collection to read
				break
			}
			return err
		}
		if bytesRead == 0 {
			// Reached last collection to read
			break
		}
		if bytesRead != len(buf) {
			return fmt.Errorf("unexpected EOF while reading name length")
		}
		nameLen, err := numbers.BytesToUint(buf)
		if err != nil {
			return err
		}

		// Read the collection name
		name, err := mustRead(nameLen, "collection name")
		if err != nil {
			return err
		}
		nameStr := string(name)

		for {
			documentDataLen, err := mustReadUint(64, "document data length")
			if err != nil {
				return err
			}

			buf, err := mustRead(documentDataLen, "document data")
			if err != nil {
				return err
			}
			bufAsBson := bson.Raw(buf)

			validationErr := bufAsBson.Validate()
			if restoreMethod != nil {
				restoreMethod(nameStr, bufAsBson, validationErr)
			} else if validationErr != nil {
				return validationErr
			}

			isLastDocument, err := mustRead(1, "is last document")
			if err != nil {
				return err
			}
			if isLastDocument[0] == 1 {
				break
			}
		}
	}

	return nil
}

// CreateBackupFile creates a backup file from the database contents
//
// !!WHAT YOU NEED TO DO!!:
//  - Close the returned file even on error
//  - Remove the returned file even on error
func CreateBackupFile(genericConn db.Connection, masterKey string) (*os.File, error) {
	backupFile, err := os.Create("./backup.gz.aes")
	if err != nil {
		return backupFile, err
	}

	conn, ok := genericConn.(*Connection)
	if !ok {
		return backupFile, errors.New("DB Connection is not a Mongo DB connection")
	}

	names, err := conn.db.ListCollectionNames(dbHelpers.Ctx(), bson.M{})
	if err != nil {
		return backupFile, err
	}

	encryptionWriter, err := crypto.NewEncryptWriter([]byte(masterKey), backupFile)
	if err != nil {
		return backupFile, err
	}

	zw, err := gzip.NewWriterLevel(encryptionWriter, 5)
	if err != nil {
		encryptionWriter.Close()
		return backupFile, err
	}

	for _, name := range names {
		ctx := dbHelpers.Ctx()
		cursor, err := conn.db.Collection(name).Find(ctx, bson.M{})
		if err != nil {
			zw.Close()
			encryptionWriter.Close()
			return backupFile, err
		}
		first := true
		for cursor.Next(ctx) {
			document := make(bson.Raw, len(cursor.Current))
			copy(document, cursor.Current)
			if first {
				// Only write the name of the collection once we are sure
				collectionNameData := append(numbers.UintToBytes(uint64(len(name)), 16), []byte(name)...)
				zw.Write(collectionNameData)

				first = false
			} else {
				// Write the is last document byte
				// In this case there is a next document so we write a false / 0
				zw.Write([]byte{0})
			}

			zw.Write(numbers.UintToBytes(uint64(len(document)), 64))
			zw.Write(document)
		}
		if !first {
			// Write the last document byte
			// Only write the last document identifier if there where actually documents in this collection
			zw.Write([]byte{1})
		}
		cursor.Close(ctx)
	}

	err1 := zw.Close()
	err2 := encryptionWriter.Close()
	if err1 != nil {
		return backupFile, err1
	}
	if err2 != nil {
		return backupFile, err2
	}

	err = backupFile.Sync()
	if err != nil {
		return backupFile, err
	}

	_, err = backupFile.Seek(0, 0)
	if err != nil {
		return backupFile, err
	}

	log.Info("validating generated backup file")

	// Validate the generated data
	err = readbackup(backupFile, masterKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate backup data: %s", err)
	}

	return backupFile, nil
}
