package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"leg3nd-pillar/internal/config"
	"leg3nd-pillar/internal/controller"
	"log"
)

func Routes(app *fiber.App) {
	public(app)

	// Assign JWT Middleware
	jwtSecret, err := config.GetEnv("JWT_SECRET")
	if err != nil {
		log.Fatalln("error occurred while parsing jwt secret env", err)
	}
	app.Use(jwtware.New(jwtware.Config{SigningKey: []byte(*jwtSecret)}))

	restricted(app)
}

func public(app *fiber.App) {
	app.Get("/ping", controller.Pong)
	app.Post("/v1/login/google", controller.LoginWithGoogle)

}

func restricted(app *fiber.App) {
	app.Post("/v1/signup", controller.CompleteSignUp)
	app.Get("/v1/me", controller.GetMyAccountInfo)
}
