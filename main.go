package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/script-development/RT-CV/controller"
	"github.com/script-development/RT-CV/db"
)

func main() {
	// Loadin the .env if available
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file:", err.Error())
	}

	// Connect to the database using the env variables
	db.ConnectToDB()

	// Create a new fiber instance (http server)
	// do not use fiber Prefork!, this app is not written to support it
	app := fiber.New()

	// Setup the app routes
	controller.Routes(app)

	// Start the webserver
	log.Fatal(app.Listen(":3000"))
}
