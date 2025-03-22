package main

import (
	"log"

	"samey/config"
	"samey/handlers"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Load .env variables
	config.LoadEnv()

	app := fiber.New()

	app.Post("/add-user", handlers.AddUser)

	log.Fatal(app.Listen(":3000"))
}
