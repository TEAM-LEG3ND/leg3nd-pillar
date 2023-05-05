package client

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"leg3nd-pillar/internal/config"
	"leg3nd-pillar/internal/dto"
	"log"
	"strconv"
)

func CreateAccount(newAccountRequest *dto.CreateAccountRequest) (*int64, error) {
	accountHost, err := config.GetEnv("ACCOUNT_HOST")
	if err != nil {
		return nil, err
	}
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(*accountHost + "/v1")
	a.JSON(newAccountRequest)
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *dto.CreateAccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("CreateAccount failed: %v", errs)
		return nil, err
	}

	log.Printf("CreateAccount: received : %v %v", statusCode, string(resultBody))

	return data.Id, nil
}

func UpdateAccount(id int64, updateAccountRequestBody *dto.UpdateAccountRequestBody) (*int64, error) {
	accountHost, err := config.GetEnv("ACCOUNT_HOST")
	if err != nil {
		return nil, err
	}
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPatch)
	req.SetRequestURI(*accountHost + "/v1/" + strconv.FormatInt(id, 10))
	a.JSON(updateAccountRequestBody)
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *dto.UpdateAccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("UpdateAccount failed: %v", errs)
		return nil, err
	}

	log.Printf("UpdateAccount: received : %v %v", statusCode, string(resultBody))

	return data.Id, nil
}

func FindAccountById(id int64) (*dto.AccountResponse, error) {
	accountHost, err := config.GetEnv("ACCOUNT_HOST")
	if err != nil {
		return nil, err
	}
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(*accountHost + "/v1/" + strconv.FormatInt(id, 10))
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *dto.AccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("FindAccountById failed: %v", errs)
		return nil, err
	}

	log.Printf("FindAccountById: received : %v %v", statusCode, string(resultBody))

	return data, nil
}

func FindAccountByEmail(email string) (*dto.AccountResponse, error) {
	accountHost, err := config.GetEnv("ACCOUNT_HOST")
	if err != nil {
		return nil, err
	}
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(*accountHost + "/v1/email/" + email)
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *dto.AccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("FindAccountByEmail failed: %v", errs)
		return nil, err
	}

	log.Printf("FindAccountByEmail: received : %v %v", statusCode, string(resultBody))

	return data, nil
}
