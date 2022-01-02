package main

import (
	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	//output, _ := helpers.ExecuteCommand("ls", 3, "-la")

	//log.Println(output)

	database.InitDB()

	app := fiber.New()
	routes.InitRoutes(app)

	app.Listen(":1337")
}
