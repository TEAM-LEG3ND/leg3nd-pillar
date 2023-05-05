package service

import (
	"fmt"
	"leg3nd-pillar/internal/client"
	"leg3nd-pillar/internal/dto"
	"log"
)

func CreateDraftAccount(email string, fullName string) (*int64, error) {
	createAccountRequest := &dto.CreateAccountRequest{
		Email:         email,
		FullName:      fullName,
		OAuthProvider: dto.AccountOAuthProviderGoogle,
	}
	createdAccountId, err := client.CreateAccount(createAccountRequest)
	if err != nil {
		log.Printf("error occurred in account service CreateDraftAccount method: %v", err)
		return nil, fmt.Errorf("error occurred in account service CreateDraftAccount method: %v", err)
	}
	return createdAccountId, nil
}

func CompleteAccountSignUp(accountId int64, nickname string, status dto.AccountStatus) (*int64, error) {
	_, err := GetAccountById(accountId)
	if err != nil {
		log.Printf("error occurred in account service CompleteAccountSignUp: %v", err)
		return nil, fmt.Errorf("error occurred in account service CompleteAccountSignUp: %v", err)
	}

	updateAccountRequestBody := &dto.UpdateAccountRequestBody{
		Nickname: &nickname,
		Status:   &status,
	}
	updatedAccountId, err := client.UpdateAccount(accountId, updateAccountRequestBody)
	if err != nil {
		log.Printf("error occurred in account service CompleteAccountSignUp while trying to update account: %v", err)
		return nil, fmt.Errorf("error occurred in account service CompleteAccountSignUp while trying to update account: %v", err)
	}

	return updatedAccountId, nil
}

func GetAccountById(accountId int64) (*dto.AccountResponse, error) {
	return client.FindAccountById(accountId)
}

func GetAccountByEmail(email string) (*dto.AccountResponse, error) {
	return client.FindAccountByEmail(email)
}
