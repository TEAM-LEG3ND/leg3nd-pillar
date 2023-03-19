package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"leg3nd-pillar/handler"
)

func Routes(app *fiber.App, corsConfig fiber.Handler) {
	api := app.Group("/api", logger.New())
	api.Use(corsConfig)
	api.Get("", handler.Auth)
	api.Get("/ping", handler.Pong)
	api.Post("/google", handler.CallbackJson)
	api.Get("/google/callback", handler.Callback)
}
