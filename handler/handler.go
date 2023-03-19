package handler

import (
	"github.com/gofiber/fiber/v2"
	"leg3nd-pillar/auth"
)

func Auth(ctx *fiber.Ctx) error {
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return ctx.Redirect(url)
}

func Callback(ctx *fiber.Ctx) error {
	token, err := auth.GetGoogleOAuthToken(ctx, ctx.FormValue("code"))
	if err != nil {
		panic(err)
	}
	user, err := auth.GetGoogleOAuthUser(token)
	if err != nil {
		panic(err)
	}
	return ctx.Status(200).JSON(user)
}
