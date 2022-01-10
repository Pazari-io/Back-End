package main

import (
	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// initalize everything
	database.InitDB()
	app := fiber.New()

	routes.InitRoutes(app)
	app.Listen("0.0.0.0:" + internal.GetKey("PORT"))
}
