package dto

type AccountResponse struct {
	Id            *int64                `json:"id"`
	Email         *string               `json:"email"`
	Nickname      *string               `json:"nickname"`
	FullName      *string               `json:"full_name"`
	OAuthProvider *AccountOAuthProvider `json:"o_auth_provider"`
	Status        *AccountStatus        `json:"status"`
}
