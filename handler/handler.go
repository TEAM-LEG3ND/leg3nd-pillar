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
		log.Printf("first error on FindAccountByEmail but do not matter maybe: %v", err)
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

func Login(ctx *fiber.Ctx) error {
	var req *model.GoogleOAuthUserRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("parse body error occurred,", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.FindAccountByEmailFailedResponse{
			Code:    model.LoginFailed,
			Message: "Parse GoogleOAuth body error occurred",
		})
	}
	token, err := auth.GetGoogleOAuthToken(ctx, req.Code)
	if err != nil {
		log.Println("GetGoogleOAuthToken error occurred,", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.FindAccountByEmailFailedResponse{
			Code:    model.LoginFailed,
			Message: "Get Google OAuth Token with Code Failed",
		})
	}
	user, err := auth.GetGoogleOAuthUser(token)
	if err != nil {
		log.Println("GetGoogleOAuthUser error occurred,", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.FindAccountByEmailFailedResponse{
			Code:    model.LoginFailed,
			Message: "Get Google User with OAuth Token Failed",
		})
	}

	ac, err := auth.FindAccountByEmail(user)
	if err != nil {
		log.Println("Cannot find account by email", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.FindAccountByEmailFailedResponse{
			Code:             model.NewUser,
			Message:          "Cannot find account by email",
			OAuthAccessToken: token.AccessToken,
			OAuthProvider:    "google",
		})
	}

	return ctx.Status(200).JSON(ac)
}

func CreateAccount(ctx *fiber.Ctx) error {
	var req *model.NewAccountRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Request. Check your request body"}
	}
	id, err := auth.CreateAccountV1(req)
	if err != nil {
		log.Printf("error occurred when creating account not existed: %v", err)
		return fmt.Errorf("error occurred when creating account not existed: %v", err)
	}
	ac, err := auth.FindAccountById(*id)
	if err != nil {
		log.Printf("error occurred when find account by id after creation: %v", err)
		return fmt.Errorf("error occurred when find account by id after creation: %v", err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(ac)
}

func Pong(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("pong")
}
