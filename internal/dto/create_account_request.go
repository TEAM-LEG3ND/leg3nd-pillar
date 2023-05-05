package dto

type CreateAccountRequest struct {
	Email         string               `json:"email"`
	FullName      string               `json:"full_name"`
	OAuthProvider AccountOAuthProvider `json:"o_auth_provider"`
}
