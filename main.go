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
	"github.com/script-development/RT-CV/helpers/random"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
)

func main() {
	// Seed the random package so generated values are "actually" random
	random.Seed()

	// Generate a server random seed used for auth
	serverSeed := random.StringBytes(64)

	// Loading the .env if available
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	var dbConn db.Connection
	useTestingDB := os.Getenv("USE_TESTING_DB")
	if useTestingDB == "true" || useTestingDB == "TRUE" {
		dbConn = mock.NewMockDB()
	} else {
		dbConn = mongo.ConnectToDB()
	}

	dbConn.RegisterEntries(
		&models.APIKey{},
		&models.Profile{},
		&models.Secret{},
	)

	models.CheckNeedToCreateSystemKeys(dbConn)

	// Create a new fiber instance (http server)
	// do not use fiber Prefork!, this app is not written to support it
	app := fiber.New(fiber.Config{
		ErrorHandler: controller.FiberErrorHandler,
	})
	app.Use(logger.New())

	// Setup the app routes
	controller.Routes(app, dbConn, serverSeed)

	// Start the webserver
	log.Fatal(app.Listen(":4000").Error())
}
