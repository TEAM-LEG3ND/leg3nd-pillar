package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"leg3nd-pillar/config"
	"leg3nd-pillar/handler"
)

func Routes(app *fiber.App) {
	public(app)

	// Assign JWT Middleware
	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte(config.Config("JWT_SECRET"))}))

	restricted(app)
}

func public(app *fiber.App) {
	app.Get("", handler.Auth)
	app.Get("/ping", handler.Pong)
	app.Post("/v1/login/google", handler.Login)

}

func restricted(app *fiber.App) {
	app.Patch("/v1/account", handler.UpdateAccount)
}
