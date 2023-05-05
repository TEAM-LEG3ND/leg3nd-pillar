package dto

type AccountStatus string

const (
	AccountStatusDraft     = AccountStatus("draft")
	AccountStatusOk        = AccountStatus("ok")
	AccountStatusSuspended = AccountStatus("suspended")
	AccountStatusWithdraw  = AccountStatus("withdraw")
)

func (s AccountStatus) String() string {
	return string(s)
}
