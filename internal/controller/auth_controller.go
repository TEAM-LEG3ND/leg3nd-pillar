package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"leg3nd-pillar/internal/config"
	"leg3nd-pillar/internal/dto"
	"leg3nd-pillar/internal/service"
	"log"
	"strconv"
	"time"
)

func LoginWithGoogle(ctx *fiber.Ctx) error {
	var req *dto.GoogleLoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		log.Println("LoginWithGoogle error occurred while body parsing", err)
		errorMessage := "Bad Request"
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.LoginErrorResponse{
			Code:    dto.ErrorCodeLoginFailed,
			Message: &errorMessage,
			Token:   nil,
		})
	}
	return service.LoginWithGoogle(ctx, req.Code)
}

func CompleteSignUp(ctx *fiber.Ctx) error {
	draftUserToken := ctx.Locals("user").(*jwt.Token) // Parsed by middleware
	claims := draftUserToken.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		message := "error occurred while parsing sub string to int"
		log.Println(message)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": message,
		})
	}
	var updateAccountRequestBody *dto.CompleteSignUpRequest
	if err := ctx.BodyParser(&updateAccountRequestBody); err != nil {
		message := "Invalid update account body"
		log.Println(message, err)
		return ctx.Status(400).JSON(fiber.Map{"message": message})
	}

	return service.CompleteSignUp(ctx, id, updateAccountRequestBody.Nickname)
}

func GetMyAccountInfo(ctx *fiber.Ctx) error {
	draftUserToken := ctx.Locals("user").(*jwt.Token) // Parsed by middleware
	claims := draftUserToken.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		message := "error occurred while parsing sub string to int"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": message,
		})
	}
	accountById, err := service.GetAccountById(id)
	if err != nil {
		message := "error occurred while getting account by id"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": message,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(*accountById)
}

func CheckToken(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		message := "error occurred while parsing sub string to int"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": message,
		})
	}
	marshaledJson, err := json.Marshal(dto.CheckTokenResponse{AccountId: id})
	if err != nil {
		message := "error occurred while marshalling check token response dto"
		log.Println(message, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": message,
		})
	}
	ctx.Set("X-Request-Account", string(marshaledJson))
	// TODO: Need to set new access token with refresh token when token is expired
	ctx.Set("Access-Token", "new token when access token is expired and refresh token is provided")
	return ctx.SendStatus(fiber.StatusOK)
}

func RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		message := "no refresh token provided with cookie"
		log.Println(message)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": message,
		})
	}
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			message := "error occurred while verifying refresh token"
			return nil, fmt.Errorf(message)
		}
		jwtSecret, err := config.GetEnv("JWT_SECRET")
		if err != nil {
			log.Fatalln("error occurred while parsing jwt secret env", err)
			return nil, err
		}
		return []byte(*jwtSecret), nil
	})
	if err != nil {
		message := "error occurred while parsing refresh token"
		log.Println(message, err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": message,
		})
	}
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		sub := claims["sub"].(string)
		id, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			message := "error occurred while parsing sub string to int"
			log.Println(message, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": message,
			})
		}
		accountById, err := service.GetAccountById(id)
		if err != nil {
			message := "error occurred while getting account by id"
			log.Println(message, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": message,
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
		jwtRefreshExpiresMinute, err := strconv.ParseInt(*jwtRefreshExpiresMinuteStr, 10, 64)
		if err != nil {
			message := "error occurred while parsing JWT refresh expires"
			log.Println(message, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(dto.LoginErrorResponse{
				Code:    dto.ErrorCodeLoginFailed,
				Message: &message,
			})
		}
		refreshToken, err := service.GetRefreshToken(*accountById.Id, time.Minute*time.Duration(jwtRefreshExpiresMinute))
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
		return ctx.SendStatus(fiber.StatusOK)
	} else {
		message := "error occurred on token validation"
		log.Println(message)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": message,
		})
	}
}
