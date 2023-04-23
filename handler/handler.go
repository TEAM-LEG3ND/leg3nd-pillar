package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"leg3nd-pillar/auth"
	"leg3nd-pillar/config"
	"leg3nd-pillar/model"
	"log"
	"strconv"
	"time"
)

func Auth(ctx *fiber.Ctx) error {
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return ctx.Redirect(url)
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

		newAccountRequest := &model.NewAccountRequest{
			Email:         user.Email,
			FullName:      user.GivenName,
			OAuthProvider: "google",
		}

		createdAccountId, err := auth.CreateAccount(newAccountRequest)

		accessToken, err := auth.GetAccessToken(*createdAccountId, time.Minute*30)
		if err != nil {
			log.Printf("error occurred on creating access token: %v", err)
			return fmt.Errorf("error occurred on creating access token: %v", err)
		}

		tokenResponse := &model.TokenResponse{AccessToken: accessToken}

		return ctx.Status(fiber.StatusUnauthorized).JSON(model.FindAccountByEmailFailedResponse{
			Code:    model.NewUser,
			Message: "Cannot find account by email",
			Token:   *tokenResponse,
		})
	}

	jwtExpiresMinute, err := strconv.ParseInt(config.Config("JWT_EXPIRES_MINUTE"), 10, 64)
	if err != nil {
		log.Printf("error occurred while parsing JWT expires")
		return fmt.Errorf("error occurred while parsing JWT expires")
	}
	jwtRefreshExpiresMinute, err := strconv.ParseInt(config.Config("JWT_REFRESH_EXPIRES_MINUTE"), 10, 64)
	if err != nil {
		log.Printf("error occurred while parsing JWT refresh expires")
		return fmt.Errorf("error occurred while parsing JWT refresh expires")
	}

	accessToken, err := auth.GetAccessToken(ac.Id, time.Minute*time.Duration(jwtExpiresMinute))
	if err != nil {
		log.Printf("error occurred on creating access token: %v", err)
		return fmt.Errorf("error occurred on creating access token: %v", err)
	}
	refreshToken, err := auth.GetRefreshToken(ac.Id, time.Minute*time.Duration(jwtRefreshExpiresMinute))
	if err != nil {
		log.Printf("error occurred on creating refresh token: %v", err)
		return fmt.Errorf("error occurred on creating refresh token: %v", err)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: *refreshToken,
		//Path:        "",
		//Domain:      "",
		MaxAge:      60 * 60 * 24 * 30,
		Expires:     time.Now().Add(time.Hour * 24 * 30),
		Secure:      true,
		HTTPOnly:    true,
		SameSite:    "none",
		SessionOnly: false,
	})

	tokenResponse := &model.TokenResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(200).JSON(tokenResponse)
}

func UpdateAccount(ctx *fiber.Ctx) error {
	draftUser := ctx.Locals("user").(*jwt.Token)
	claims := draftUser.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		log.Printf("error occurred when parsing jwt sub into int64, %v", err)
		return fmt.Errorf("error occurred when parsing jwt sub into int64, %w", err)
	}

	_, err = auth.FindAccountById(id)
	if err != nil {
		log.Printf("error occurred when find account by id after creation: %v", err)
		return fmt.Errorf("error occurred when find account by id after creation: %v", err)
	}

	var updateAccountRequestBody *model.UpdateAccountRequestBody
	if err := ctx.BodyParser(&updateAccountRequestBody); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Invalid update account body"})
	}

	_, err = auth.UpdateAccount(id, updateAccountRequestBody)
	if err != nil {
		log.Printf("error occurred while updating account by id: %v", err)
		return fmt.Errorf("error occurred while updating account by id: %v", err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func Pong(ctx *fiber.Ctx) error {
	return ctx.Status(200).SendString("pong")
}
