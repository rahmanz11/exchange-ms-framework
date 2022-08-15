package models

// Payload to handle requests containing the account_number
type AccountNumberPayload struct {
	AccountNumber string `json:"account_number"`
}

// Payload to handle requests containing the transaction_id
type TransactionIdPayload struct {
	TransactionId string `json:"transaction_id"`
}

// Payload to handle requests for Linking Accounts
type LinkedAccountPayload struct {
	AccountNumber    string `json:"account_number"`
	SubAccountNumber string `json:"sub_account_number"`
}
