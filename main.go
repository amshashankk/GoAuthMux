package main

import (
	"github.com/amshashankk/database"
	"github.com/amshashankk/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	client := database.Connect("MyDatabase", "mongodb://localhost:27017")
	if client == nil {
		return
	}
	//defer database.Disconnect(client) // Disconnecting once the main finished execution

	app := fiber.New()
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	err := app.Listen(":8000")
	if err != nil {
		return
	}
}
