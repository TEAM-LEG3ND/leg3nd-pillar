package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"leg3nd-pillar/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	corsConfig := cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
	app.Use(corsConfig)
	router.Routes(app, corsConfig)
	app.Listen(":8081")
}
