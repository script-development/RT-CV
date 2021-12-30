package backup

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/mongo"
	"github.com/script-development/RT-CV/models"
)

// StartScheduleOptions are the options required to run StartsSchedule
// All fields are required to start the backup schedule
type StartScheduleOptions struct {
	// The key used to encrypt / decrypt the backup files
	BackupEncryptionKey string

	// S3 connection options
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Bucket          string
	S3UseSSL          bool
}

// StartScheduleOptionsFromEnv creates the StartScheduleOptions from the environment variables
func StartScheduleOptionsFromEnv() StartScheduleOptions {
	return StartScheduleOptions{
		BackupEncryptionKey: os.Getenv("MONGODB_BACKUP_KEY"),
		S3Endpoint:          os.Getenv("BACKUP_S3_ENDPOINT"),
		S3AccessKeyID:       os.Getenv("BACKUP_S3_ACCESS_KEY_ID"),
		S3SecretAccessKey:   os.Getenv("BACKUP_S3_SECRET_ACCESS_KEY"),
		S3Bucket:            os.Getenv("BACKUP_S3_BUCKET"),
		S3UseSSL:            strings.ToLower(os.Getenv("BACKUP_S3_USE_SSL")) == "true",
	}
}

func (o *StartScheduleOptions) createS3Client(createBucketIfNotExist bool) *minio.Client {
	s3Client, err := minio.New(o.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(o.S3AccessKeyID, o.S3SecretAccessKey, ""),
		Secure: o.S3UseSSL,
	})
	if err != nil {
		log.WithError(err).Fatal("failed to create minio client")
	}

	if createBucketIfNotExist {
		bucketExists, err := s3Client.BucketExists(context.Background(), o.S3Bucket)
		if err != nil {
			log.WithError(err).Fatal("failed to check if backup bucket exists")
		}

		if !bucketExists {
			// Try to create the bucket if it doesn't exist yet
			err = s3Client.MakeBucket(context.Background(), o.S3Bucket, minio.MakeBucketOptions{})
			if err != nil {
				log.WithError(err).Fatal("unable to create the bucket used to store backups")
			}
		}
	}

	return s3Client
}

func unwrapMongoConn(dbConn db.Connection) *mongo.Connection {
	conn, ok := dbConn.(*mongo.Connection)
	if !ok {
		log.Fatal("DB Connection is not a Mongo DB connection")
	}
	return conn
}

// StartsSchedule starts the cron job for creating the backups
// The backupMasterKey is used to encrypt the backup files generated
func StartsSchedule(dbConn db.Connection, options StartScheduleOptions, forceBackup bool) {
	mongoDBConn := unwrapMongoConn(dbConn)

	if len(options.BackupEncryptionKey) < 16 {
		msg := "encryption key is too short, make sure you have set the MONGODB_BACKUP_KEY env variable"
		log.Fatalf("Error initializing backup: " + msg)
	}

	s3Client := options.createS3Client(true)

	// Check every 24 hours if we need to create a backup
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		checkNeedBackup(s3Client, options.S3Bucket, mongoDBConn, options.BackupEncryptionKey, forceBackup)
		for range ticker.C {
			checkNeedBackup(s3Client, options.S3Bucket, mongoDBConn, options.BackupEncryptionKey, false)
		}
	}()
}

func checkNeedBackup(s3Client *minio.Client, bucketName string, dbConn *mongo.Connection, backupMasterKey string, force bool) {
	if force {
		log.Info("creating a new backup..")
	} else {
		needToCreateDB, err := models.NeedToCreateBackup(dbConn)
		if err != nil {
			log.WithError(err).Error("Failed to check if backup is needed")
			return
		}
		if !needToCreateDB {
			return
		}

		log.Info("to long ago since last backup, creating a new backup..")
	}

	backupFile, err := CreateBackupFile(dbConn, backupMasterKey)
	if err != nil {
		log.WithError(err).Error("Failed to create backup of database")
		return
	}
	defer func() {
		backupFile.Close()
		os.Remove(backupFile.Name())
	}()

	backupFileStat, err := backupFile.Stat()
	if err != nil {
		log.WithError(err).Error("Failed to get meta information about the created backup file")
		return
	}

	bucketFileName := fmt.Sprintf(
		"/rt-cv-backups/%s.gz.aes",
		time.Now().Format("2006-01-02--15-04"),
	)

	_, err = s3Client.PutObject(
		context.Background(),
		bucketName,
		bucketFileName,
		backupFile,
		backupFileStat.Size(),
		minio.PutObjectOptions{
			ContentType: "binary/octet-stream",
		},
	)
	if err != nil {
		log.WithError(err).Error("Failed to upload backup file to S3")
	} else {
		log.Infof("uploaded backup file to S3 with name %s", bucketFileName)
	}
}
