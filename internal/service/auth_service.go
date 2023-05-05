package service

import (
	"github.com/gofiber/fiber/v2"
	"leg3nd-pillar/internal/client"
	"leg3nd-pillar/internal/config"
	"leg3nd-pillar/internal/dto"
	"log"
	"strconv"
	"time"
)

func LoginWithGoogle(ctx *fiber.Ctx, code string) error {
	token, err := client.GetGoogleOAuthToken(ctx, code)
	if err != nil {
		message := "Get Google OAuth Token with Code Failed"
		log.Println(message, err)

		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	user, err := client.GetGoogleOAuthUser(token)
	if err != nil {
		message := "Get Google User with OAuth Token Failed"
		log.Println(message, err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}

	account, err := client.FindAccountByEmail(user.Email)
	if err != nil {
		message := "Cannot find account by email"
		log.Println(message, err)

		draftAccountId, err := CreateDraftAccount(user.Email, user.Name)
		if err != nil {
			message := "Create Draft Account failed"
			log.Println(message, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
				Code:    dto.ErrorCodeLoginFailed,
				Message: &message,
			})
		}

		accessToken, err := GetAccessToken(*draftAccountId, time.Minute*30)
		if err != nil {
			message := "error occurred while creating access token"
			log.Println(message, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
				Code:    dto.ErrorCodeLoginFailed,
				Message: &message,
			})
		}

		tokenResponse := dto.TokenResponse{AccessToken: accessToken}

		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeNewUser,
			Message: &message,
			Token:   &tokenResponse,
		})
	} else if *account.Status == dto.AccountStatusDraft {
		message := "Account is found but in draft status"
		accessToken, err := GetAccessToken(*account.Id, time.Minute*30)
		if err != nil {
			message := "error occurred while creating access token"
			log.Println(message, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
				Code:    dto.ErrorCodeLoginFailed,
				Message: &message,
			})
		}

		tokenResponse := dto.TokenResponse{AccessToken: accessToken}

		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeNewUser,
			Message: &message,
			Token:   &tokenResponse,
		})
	}

	jwtExpiresMinuteStr, err := config.GetEnv("JWT_EXPIRES_MINUTE")
	if err != nil {
		message := "error occurred while parsing JWT expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	jwtRefreshExpiresMinuteStr, err := config.GetEnv("JWT_REFRESH_EXPIRES_MINUTE")
	if err != nil {
		message := "error occurred while parsing JWT refresh expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}

	jwtExpiresMinute, err := strconv.ParseInt(*jwtExpiresMinuteStr, 10, 64)
	if err != nil {
		message := "error occurred while parsing JWT expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	jwtRefreshExpiresMinute, err := strconv.ParseInt(*jwtRefreshExpiresMinuteStr, 10, 64)
	if err != nil {
		message := "error occurred while parsing JWT refresh expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}

	accessToken, err := GetAccessToken(*account.Id, time.Minute*time.Duration(jwtExpiresMinute))
	if err != nil {
		message := "error occurred while creating jwt token"
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	refreshToken, err := GetRefreshToken(*account.Id, time.Minute*time.Duration(jwtRefreshExpiresMinute))
	if err != nil {
		message := "error occurred while creating jwt refresh token"
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
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

	tokenResponse := &dto.TokenResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(fiber.StatusOK).JSON(tokenResponse)
}

func CompleteSignUp(ctx *fiber.Ctx, accountId int64, nickname string) error {
	updatedAccountId, err := CompleteAccountSignUp(accountId, nickname, dto.AccountStatusOk)
	if err != nil {
		message := "error occurred while updating account"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": message})
	}

	jwtExpiresMinuteStr, err := config.GetEnv("JWT_EXPIRES_MINUTE")
	if err != nil {
		message := "error occurred while parsing JWT expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	jwtRefreshExpiresMinuteStr, err := config.GetEnv("JWT_REFRESH_EXPIRES_MINUTE")
	if err != nil {
		message := "error occurred while parsing JWT refresh expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}

	jwtExpiresMinute, err := strconv.ParseInt(*jwtExpiresMinuteStr, 10, 64)
	if err != nil {
		message := "error occurred while parsing JWT expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	jwtRefreshExpiresMinute, err := strconv.ParseInt(*jwtRefreshExpiresMinuteStr, 10, 64)
	if err != nil {
		message := "error occurred while parsing JWT refresh expires"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}

	accessToken, err := GetAccessToken(*updatedAccountId, time.Minute*time.Duration(jwtExpiresMinute))
	if err != nil {
		message := "error occurred while creating jwt token"
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
	}
	refreshToken, err := GetRefreshToken(*updatedAccountId, time.Minute*time.Duration(jwtRefreshExpiresMinute))
	if err != nil {
		message := "error occurred while creating jwt refresh token"
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &message,
		})
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

	tokenResponse := &dto.TokenResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(fiber.StatusOK).JSON(tokenResponse)
}
