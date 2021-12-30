package backup

import (
	"context"
	"path"

	"github.com/apex/log"
	"github.com/minio/minio-go/v7"
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)

// Restore restores a s3 backup to the mongodb database
func Restore(dbConn db.Connection, backupFile string, options StartScheduleOptions) {
	mongoDB := unwrapMongoConn(dbConn).GetDB()
	s3Client := options.createS3Client(true)

	log.Infof("Restoring backup %s from s3..", backupFile)

	_, err := s3Client.StatObject(context.Background(), options.S3Bucket, backupFile, minio.StatObjectOptions{})
	if err != nil {
		backupFile = path.Join("/rt-cv-backups/", backupFile)
		_, err2 := s3Client.StatObject(context.Background(), options.S3Bucket, backupFile, minio.StatObjectOptions{})
		if err2 != nil {
			log.WithError(err).Fatal("unable to get backup file from s3")
		}
	}

	obj, err := s3Client.GetObject(context.Background(), options.S3Bucket, backupFile, minio.GetObjectOptions{})
	if err != nil {
		log.WithError(err).Fatal("failed to get the backup file from s3")
	}

	log.Infof("Found backup file in s3, verifying it's content..")

	err = readbackup(obj, options.BackupEncryptionKey, nil)
	if err != nil {
		log.WithError(err).Fatal("failed to validate the backup file")
	}

	log.Infof("Content of backup file is valid, inserting data into temp collection..")

	collectionNamesList, err := mongoDB.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.WithError(err).Fatal("unable to get the list of collections currently in the database")
	}
	hasCollection := func(name string) bool {
		for _, collectionName := range collectionNamesList {
			if collectionName == name {
				return true
			}
		}
		return false
	}

	restoreCollections := map[string]*mongodb.Collection{}
	removeTempCollections := func() {
		log.Info("trying to removing temp collections..")
		for _, collection := range restoreCollections {
			collection.Drop(context.Background())
		}
	}

	obj.Seek(0, 0)
	err = readbackup(obj, options.BackupEncryptionKey, func(collectionName string, doc bson.Raw, err error) {
		collectionNameToRestoreTo := collectionName + "_RESTORE_TEMP"
		collectionToRestoreTo, ok := restoreCollections[collectionNameToRestoreTo]
		if !ok {
			if hasCollection(collectionNameToRestoreTo) {
				err := mongoDB.Collection(collectionNameToRestoreTo).Drop(context.Background())
				if err != nil {
					log.WithError(err).Fatal("unable to drop earlier temp collection")
				}
			}

			err := mongoDB.CreateCollection(context.Background(), collectionNameToRestoreTo)
			if err != nil {
				removeTempCollections()
				log.WithError(err).Fatalf("failed to create temporary collection for the %s collection", collectionNameToRestoreTo)
			}
			collectionToRestoreTo = mongoDB.Collection(collectionNameToRestoreTo)
			restoreCollections[collectionNameToRestoreTo] = collectionToRestoreTo
		}

		_, err = collectionToRestoreTo.InsertOne(context.Background(), doc)
		if err != nil {
			removeTempCollections()
			log.WithError(err).
				WithField("collection", collectionNameToRestoreTo).
				Fatal("failed to insert a document into the temporary restore collection")
		}
	})

	helperText := "If anything after this fails you can manually promote the temp collections to the " +
		"real collections using the following command (you need to rename some_collection_name to the db collection names): " +
		"db.some_collection_name_RESTORE_TEMP.renameCollection('some_collection_name', true)"

	// How to manually restore:
	// 1. Open the mongo shell/console
	// 2. Run: use rt-cv
	// 3. For every collection with the _RESTORE_TEMP suffix, run:
	//      db.some_collection_name_RESTORE_TEMP.renameCollection('some_collection_name', true)

	log.Info("Inserting temp data succeeded, doing the same on the real collections.. " + helperText)

	restoredDocs := 0

	collections := map[string]*mongodb.Collection{}

	obj.Seek(0, 0)
	err = readbackup(obj, options.BackupEncryptionKey, func(collectionName string, doc bson.Raw, err error) {
		collection, ok := collections[collectionName]
		if !ok {
			if hasCollection(collectionName) {
				err = mongoDB.Collection(collectionName).Drop(context.Background())
				if err != nil {
					log.WithError(err).Fatal("unable to drop collection so we can re-create it and restore data")
				}
			}
			err := mongoDB.CreateCollection(context.Background(), collectionName)
			if err != nil {
				log.WithError(err).Fatal("failed to create collection " + collection.Name())
			}
			collection = mongoDB.Collection(collectionName)
			collections[collectionName] = collection
		}

		restoredDocs++
		_, err = collection.InsertOne(context.Background(), doc)
		if err != nil {
			log.WithError(err).WithField("collection", collectionName).Fatal("failed to insert a document")
		}
	})

	removeTempCollections()

	log.Infof("Restore successfull, restored %d documents over %d collections", restoredDocs, len(collections))
}
