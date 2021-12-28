package backup

import (
	"time"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/models"
)

// StartsSchedule starts the cron job for creating the backups
// The backupMasterKey is used to encrypt the backup files generated
func StartsSchedule(dbConn db.Connection, backupMasterKey string) {
	if len(backupMasterKey) < 16 {
		msg := "encryption key is too short, make sure you have set the MONGODB_BACKUP_KEY env variable"
		log.Fatalf("Error initializing backup: " + msg)
	}

	// Check every 24 hours if we need to create a backup
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		checkNeedBackup(dbConn, backupMasterKey)
		for range ticker.C {
			checkNeedBackup(dbConn, backupMasterKey)
		}
	}()
}

func checkNeedBackup(dbConn db.Connection, backupMasterKey string) {
	needToCreateDB, err := models.NeedToCreateBackup(dbConn)
	if err != nil {
		log.WithError(err).Error("Failed to check if backup is needed")
		return
	}
	if !needToCreateDB {
		return
	}

	log.Info("to long ago since last backup, creating a new backup..")

	backupFile, err := CreateBackupFile(dbConn, backupMasterKey)
	if err != nil {
		log.WithError(err).Error("Failed to create backup of database")
		return
	}
	backupFile.Close()
	log.Infof("created backup file with name %s", backupFile.Name())
}
