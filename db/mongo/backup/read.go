package backup

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/script-development/RT-CV/helpers/crypto"
	"github.com/script-development/RT-CV/helpers/numbers"
	"go.mongodb.org/mongo-driver/bson"
)

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
