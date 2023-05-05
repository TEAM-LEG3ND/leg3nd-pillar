package dto

type ErrorCode string

const (
	ErrorCodeNewUser     = ErrorCode("NEW_USER")
	ErrorCodeLoginFailed = ErrorCode("LOGIN_FAILED")
)

func (ec ErrorCode) String() string {
	return string(ec)
}

type LoginErrorResponse struct {
	Code    ErrorCode      `json:"error_code"`
	Message *string        `json:"message"`
	Token   *TokenResponse `json:"token"`
}
