package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"leg3nd-pillar/handler"
)

func Routes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("", handler.Auth)
	api.Get("/google/callback", handler.Callback)
}
