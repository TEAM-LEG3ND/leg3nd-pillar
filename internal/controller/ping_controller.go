package controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func Pong(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()
	log.Println(headers)
	return ctx.Status(fiber.StatusOK).SendString("pong, " + headers["X-Account-Id"])
}
