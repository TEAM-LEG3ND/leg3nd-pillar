package dto

type AccountOAuthProvider string

const (
	AccountOAuthProviderGoogle = AccountOAuthProvider("google")
	AccountOAuthProviderGitHub = AccountOAuthProvider("github")
)

func (oap AccountOAuthProvider) String() string {
	return string(oap)
}
