package controller

import "github.com/gofiber/fiber/v2"

func Pong(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("pong")
}
