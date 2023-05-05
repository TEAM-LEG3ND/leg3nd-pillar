package client

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"leg3nd-pillar/internal/config"
	"leg3nd-pillar/internal/dto"
	"log"
)

func GetGoogleOAuthToken(ctx *fiber.Ctx, code string) (*oauth2.Token, error) {
	googleOAuthConfig, err := getGoogleOAuthConfig()
	if err != nil {
		return nil, err
	}
	token, err := googleOAuthConfig.Exchange(ctx.Context(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetGoogleOAuthUser(token *oauth2.Token) (*dto.GoogleOAuthUserResponse, error) {
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
	var data *dto.GoogleOAuthUserResponse

	if statusCode, resultBody, errs = a.Struct(&data); len(errs) > 0 {
		err := fmt.Errorf("GetGoogleOAuthUser failed: %v", errs)
		return nil, err
	}

	log.Printf("GetGoogleOAuthUser: received : %v %v", statusCode, string(resultBody))

	return data, nil
}

func getGoogleOAuthConfig() (*oauth2.Config, error) {
	googleClientId, err := config.GetEnv("GOOGLE_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	googleClientSecret, err := config.GetEnv("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	googleRedirectUrl, err := config.GetEnv("GOOGLE_REDIRECT_URL")
	if err != nil {
		return nil, err
	}
	googleScopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	conf := &oauth2.Config{
		ClientID:     *googleClientId,
		ClientSecret: *googleClientSecret,
		RedirectURL:  *googleRedirectUrl,
		Scopes:       googleScopes,
		Endpoint:     google.Endpoint,
	}

	return conf, nil
}
