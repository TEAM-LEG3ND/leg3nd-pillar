package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"leg3nd-pillar/internal/dto"
	"leg3nd-pillar/internal/service"
	"log"
	"strconv"
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
	return ctx.Status(fiber.StatusOK).SendString(strconv.FormatInt(id, 10))
}
