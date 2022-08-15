package payloads

// LinkAccountPayload to be used for the link account api routes
type LinkAccountPayload struct {
	AccountNumber    string `json:"account_number"`
	SubAccountNumber string `json:"sub_account_number"`
}
