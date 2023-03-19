package model

// GoogleResponse is the response of google oauth userinfo api
type GoogleResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

type GoogleOAuthUserRequest struct {
	Code string `json:"code"`
}
