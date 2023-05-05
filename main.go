package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"leg3nd-pillar/internal/middleware"
	"leg3nd-pillar/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Use(middleware.Cors())
	router.Routes(app)
	app.Listen(":8081")
}
