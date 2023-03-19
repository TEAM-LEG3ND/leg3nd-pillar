package handler

import (
	"github.com/gofiber/fiber/v2"
	"leg3nd-pillar/auth"
	"leg3nd-pillar/model"
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

func CallbackJson(ctx *fiber.Ctx) error {
	var req *model.GoogleOAuthUserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		panic(err)
	}
	token, err := auth.GetGoogleOAuthToken(ctx, req.Code)
	if err != nil {
		panic(err)
	}
	user, err := auth.GetGoogleOAuthUser(token)
	if err != nil {
		panic(err)
	}
	return ctx.Status(200).JSON(user)
}

func Pong(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("pong")
}
