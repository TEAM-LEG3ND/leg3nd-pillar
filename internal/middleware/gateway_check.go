package middleware

import "github.com/gofiber/fiber/v2"

func NewGatewayCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		if headers["X-From-Gateway"] != "true" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}
