package routes

import (
	"project1/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetUp(app *fiber.App) {

	app.Get("/", controllers.Hello)
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Post("/logout", controllers.Logout)
	app.Get("/user", controllers.User)

	app.Listen(":8000")
}
