package models

import "github.com/google/uuid"

type SubAccount struct {
	SubAccountId  uuid.UUID `json:"sub_account_id"`
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	Status        string    `json:"status"`
}
