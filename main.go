package main

import (
	"os"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/script-development/RT-CV/controller"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/db/mongo"
	"github.com/script-development/RT-CV/helpers/emailservice"
	"github.com/script-development/RT-CV/helpers/random"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
)

// AppVersion is used for the X-App-Version header
// This variable can be set by:
//   go build -ldflags "-X main.AppVersion=1.0.0"
var AppVersion = "LOCAL"

func main() {
	// Seed the random package so generated values are "actually" random
	random.Seed()

	// Loading the .env if available
	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err.Error())
		}
	} else {
		log.Info("No .env file found")
	}

	// Initialize the mail service
	err = emailservice.Setup(
		emailservice.EmailServerConfiguration{
			Identity: os.Getenv("EMAIL_IDENTITY"),
			Username: os.Getenv("EMAIL_USER"),
			Password: os.Getenv("EMAIL_PASSWORD"),
			Host:     os.Getenv("EMAIL_HOST"),
			Port:     os.Getenv("EMAIL_PORT"),
			From:     os.Getenv("EMAIL_FROM"),
		},
		nil,
	)
	if err != nil {
		log.WithError(err).Error("Error initializing email service")
		return
	}

	// Initialize the database
	var dbConn db.Connection
	useTestingDB := os.Getenv("USE_TESTING_DB")
	if useTestingDB == "true" || useTestingDB == "TRUE" {
		dbConn = mock.NewMockDB()
		log.WithField("id", mock.DashboardKey.ID.Hex()).WithField("key", mock.DashboardKey.Key).Info("Mock dashboard key")
	} else {
		dbConn = mongo.ConnectToDB()
	}

	dbConn.RegisterEntries(
		&models.APIKey{},
		&models.Profile{},
		&models.Secret{},
		&models.Match{},
	)

	models.CheckDashboardKeyExists(dbConn)

	// Create a new fiber instance (http server)
	// do not use fiber Prefork!, this app is not written to support it
	app := fiber.New(fiber.Config{
		ErrorHandler: controller.FiberErrorHandler,
	})
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		err = c.Next()
		c.Set("X-App-Version", AppVersion)
		return err
	})

	// Setup the app routes
	controller.Routes(app, AppVersion, dbConn, false)

	testingDieAfterInit := os.Getenv("TESTING_DIE_AFTER_INIT")
	if testingDieAfterInit == "true" || testingDieAfterInit == "TRUE" {
		return
	}

	// Start the webserver
	log.Fatal(app.Listen(":4000").Error())
}
