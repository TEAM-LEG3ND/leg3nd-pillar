package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"leg3nd-pillar/auth"
	"leg3nd-pillar/model"
	"log"
)

func Auth(ctx *fiber.Ctx) error {
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return ctx.Redirect(url)
}

func CallbackJson(ctx *fiber.Ctx) error {
	var req *model.GoogleOAuthUserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		fmt.Println(err)
	}
	token, err := auth.GetGoogleOAuthToken(ctx, req.Code)
	if err != nil {
		fmt.Println(err)
	}
	user, err := auth.GetGoogleOAuthUser(token)
	if err != nil {
		fmt.Println(err)
	}

	ac, err := auth.FindAccountByEmail(user)
	if err != nil {
		id, err := auth.CreateAccount(user)
		if err != nil {
			log.Printf("error occurred when creating account not existed: %v", err)
			return fmt.Errorf("error occurred when creating account not existed: %v", err)
		}
		ac, err = auth.FindAccountById(*id)
		if err != nil {
			log.Printf("error occurred when find account by id after creation: %v", err)
			return fmt.Errorf("error occurred when find account by id after creation: %v", err)
		}
	}

	return ctx.Status(200).JSON(ac)
}

func Pong(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("pong")
}
