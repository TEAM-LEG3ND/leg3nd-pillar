package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"leg3nd-pillar/config"
	"leg3nd-pillar/model"
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

	fmt.Println(google.Endpoint)

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

	fmt.Printf("received : %v %v", statusCode, string(resultBody))

	return data, nil
}
