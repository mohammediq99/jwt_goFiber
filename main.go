package main

import (
	"project1/database"
	"project1/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	database.Connect()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.SetUp(app)

}
