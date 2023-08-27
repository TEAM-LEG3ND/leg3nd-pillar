package controller

import "github.com/gofiber/fiber/v2"

func Pong(ctx *fiber.Ctx) error {
	xAccountId := ctx.GetReqHeaders()["x-account-id"]
	return ctx.Status(fiber.StatusOK).SendString("pong, " + xAccountId)
}
