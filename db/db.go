package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalln("Unable to parse port DB_PORT:", err.Error())
	}

	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "root"
	}

	address := os.Getenv("DB_HOST")
	if address == "" {
		address = "localhost"
	}

	if port > 0 {
		address += ":" + strconv.Itoa(port)
	} else {
		address += ":3306"
	}

	name := os.Getenv("DB_DATABASE")
	if name == "" {
		name = "first2find"
	}

	dns := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		username,
		os.Getenv("DB_PASSWORD"),
		address,
		name,
	)

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to connect to database, err:", err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalln("unable to get database reference, err:", err.Error())
	}
	sqlDB.SetMaxOpenConns(50)
}
