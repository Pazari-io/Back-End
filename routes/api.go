package routes

import (
	"github.com/Pazari-io/Back-End/handlers"
	"github.com/Pazari-io/Back-End/middlewares"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app fiber.Router) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hi 👋!")
	})

	// health check
	app.Get("/api/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Ok")
	})

	// API related routes
	baseAPI := app.Group("/api/v1")

	// Auth there is a 50 request per M limitation to avoid bots and brute forces
	authAPI := baseAPI.Group("/auth", middlewares.KeyProtected())
	authAPI.Post("/upload", handlers.Uploader)
	authAPI.Get("/watermark", handlers.DownloadWaterMarked)

}
