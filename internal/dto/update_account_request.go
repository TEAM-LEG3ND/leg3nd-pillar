package dto

type UpdateAccountRequestBody struct {
	Email         *string               `json:"email"`
	Nickname      *string               `json:"nickname"`
	FullName      *string               `json:"full_name"`
	OAuthProvider *AccountOAuthProvider `json:"o_auth_provider"`
	Status        *AccountStatus        `json:"status"`
}
