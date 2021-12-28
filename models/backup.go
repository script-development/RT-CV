package models

import (
	"time"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/mongo"
)

/*

This file contains the code for the backups collection in the database
This file mainly provides meta data for the actual backup logic in ./db/mongo/backup/create.go

*/

// Backup tells when a backup was created
type Backup struct {
	db.M `bson:",inline"`
	Time time.Time
}

// CollectionName returns the collection name of the Backup
func (*Backup) CollectionName() string {
	return "backups"
}

// LastBackupTime returns the last time of a created backup
// a zero time is returned when there is no backup found
// a zero time can be detected using the (time.Time).IsZero() method
func LastBackupTime(dbConn db.Connection) (time.Time, error) {
	backup := Backup{}
	err := dbConn.FindOne(&backup, nil)
	if err == mongo.ErrNoDocuments {
		// Return zero time
		return time.Time{}, nil
	} else if err != nil {
		// Return zero time
		return time.Time{}, err
	}
	return backup.Time, nil
}

// NeedToCreateBackup returns whether a backup is needed
func NeedToCreateBackup(dbConn db.Connection) (bool, error) {
	lastBackupDateTime, err := LastBackupTime(dbConn)
	if err != nil {
		return false, err
	}

	// It returns true if the last backup was more than a week ago
	return time.Now().AddDate(0, 0, -7).After(lastBackupDateTime), nil
}

// SetLastBackupToNow sets the last backup time to now
func SetLastBackupToNow(dbConn db.Connection) error {
	backup := Backup{}
	err := dbConn.FindOne(&backup, nil)
	if err == mongo.ErrNoDocuments {
		backup = Backup{
			M:    db.NewM(),
			Time: time.Now(),
		}
		return dbConn.Insert(&backup)
	} else if err != nil {
		return err
	}
	backup.Time = time.Now()
	return dbConn.UpdateByID(&backup)
}
