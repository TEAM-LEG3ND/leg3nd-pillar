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

type AccountResponse struct {
	Id            int64  `json:"id"`
	Email         string `json:"email"`
	Nickname      string `json:"nickname"`
	FullName      string `json:"full_name"`
	OAuthProvider string `json:"o_auth_provider"`
}

type NewAccountResponse struct {
	Id int64 `json:"id"`
}

type ErrorCode string

const (
	NewUser     = ErrorCode("NEW_USER")
	LoginFailed = ErrorCode("LOGIN_FAILED")
)

type FindAccountByEmailFailedResponse struct {
	Code             ErrorCode `json:"error_code"`
	Message          string    `json:"string"`
	OAuthAccessToken string    `json:"o_auth_access_token"`
	OAuthProvider    string    `json:"o_auth_provider"`
}
