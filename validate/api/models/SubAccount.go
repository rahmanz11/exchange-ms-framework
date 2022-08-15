package models
// SubAccount model 
type SubAccount struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Status        string  `json:"status"`
}
