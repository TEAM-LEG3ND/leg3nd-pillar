package router

import (
	"github.com/gofiber/fiber/v2"
	"leg3nd-pillar/handler"
)

func Routes(app *fiber.App) {
	app.Get("", handler.Auth)
	app.Get("/ping", handler.Pong)
	app.Post("/google", handler.CallbackJson)
}
