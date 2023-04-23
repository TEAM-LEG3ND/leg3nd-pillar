package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"leg3nd-pillar/config"
	"leg3nd-pillar/model"
	"log"
	"strconv"
	"time"
)

// ConfigGoogle returns oauth2 Config related to google from user dotenv file
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("GOOGLE_CLIENT_ID"),
		ClientSecret: config.Config("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  config.Config("GOOGLE_REDIRECT_URL"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return conf
}

func GetGoogleOAuthToken(ctx *fiber.Ctx, code string) (*oauth2.Token, error) {
	token, err := ConfigGoogle().Exchange(ctx.Context(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetGoogleOAuthUser(token *oauth2.Token) (*model.GoogleResponse, error) {
	accessToken := token.AccessToken
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI("https://www.googleapis.com/oauth2/v3/userinfo")
	req.Header.Set("Authorization", bearerToken)

	if err := a.Parse(); err != nil {
		return nil, err
	}

	var statusCode int
	var resultBody []byte
	var errs []error
	var data *model.GoogleResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("GetGoogleOAuthUser failed: %v", errs)
		return nil, err
	}

	log.Printf("GetGoogleOAuthUser: received : %v %v", statusCode, string(resultBody))

	return data, nil
}

func CreateAccount(newAccountRequest *model.NewAccountRequest) (*int64, error) {
	accountHost := config.Config("ACCOUNT_HOST")
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(accountHost + "/v1")
	a.JSON(newAccountRequest)
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *model.NewAccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("CreateAccount failed: %v", errs)
		return nil, err
	}

	log.Printf("CreateAccount: received : %v %v", statusCode, string(resultBody))

	return &data.Id, nil
}

func UpdateAccount(id int64, updateAccountRequestBody *model.UpdateAccountRequestBody) (*int64, error) {
	accountHost := config.Config("ACCOUNT_HOST")
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPatch)
	req.SetRequestURI(accountHost + "/v1/" + strconv.FormatInt(id, 10))
	a.JSON(updateAccountRequestBody)
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *model.UpdatedAccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("UpdateAccount failed: %v", errs)
		return nil, err
	}

	log.Printf("UpdateAccount: received : %v %v", statusCode, string(resultBody))

	return &data.Id, nil
}

func FindAccountById(id int64) (*model.AccountResponse, error) {
	accountHost := config.Config("ACCOUNT_HOST")
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(accountHost + "/v1/" + strconv.FormatInt(id, 10))
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *model.AccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("FindAccountById failed: %v", errs)
		return nil, err
	}

	log.Printf("FindAccountById: received : %v %v", statusCode, string(resultBody))

	return data, nil
}

func FindAccountByEmail(googleResponse *model.GoogleResponse) (*model.AccountResponse, error) {
	accountHost := config.Config("ACCOUNT_HOST")
	email := googleResponse.Email
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(accountHost + "/v1/email/" + email)
	if err := a.Parse(); err != nil {
		return nil, err
	}
	var statusCode int
	var resultBody []byte
	var errs []error
	var data *model.AccountResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("FindAccountByEmail failed: %v", errs)
		return nil, err
	}

	log.Printf("FindAccountByEmail: received : %v %v", statusCode, string(resultBody))

	return data, nil
}

func GetAccessToken(id int64, duration time.Duration) (*string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(id, 10),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config("JWT_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("token generation failed, %w", err)
	}
	return &t, nil
}

func GetRefreshToken(id int64, duration time.Duration) (*string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(id, 10),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("token generation failed, %w", err)
	}
	return &t, nil
}
