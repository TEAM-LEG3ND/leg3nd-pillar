package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"leg3nd-pillar/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	router.Routes(app)
	app.Listen(":8080")
}
