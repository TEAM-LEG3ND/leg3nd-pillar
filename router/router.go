package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"leg3nd-pillar/handler"
	"leg3nd-pillar/middleware"
)

func Routes(app *fiber.App) {
	api := app.Group("/api", logger.New(), middleware.Cors())
	api.Get("", handler.Auth)
	api.Get("/ping", handler.Pong)
	api.Post("/google", handler.CallbackJson)
	api.Get("/google/callback", handler.Callback)
}
