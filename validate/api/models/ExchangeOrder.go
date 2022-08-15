package models

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// ExchangeOrder model
type ExchangeOrder struct {
	TransactionId uuid.UUID `json:"transaction_id"`
	From          string    `json:"from"`
	Fund          string    `json:"fund"`
	Amount        float64   `json:"amt"`
	Re            string    `json:"re"`
	ReceivedAt    time.Time `json:"receivedAt"`
	ValidatedAt   time.Time `json:"validatedAt"`
}

// ValidateSubAccountStatus validates the status of the sub account to be ACTIVE
func ValidateSubAccountStatus(accountNumber string) (SubAccount, error) {
	// Get the sub account from the sub account api
	subAccountRaw, err := http.Get(os.Getenv("SUB_ACCOUNT_API") + "/api/v1/" + accountNumber)
	if err != nil {
		return SubAccount{}, err
	}
	// Read the body of the response
	body, err := ioutil.ReadAll(subAccountRaw.Body)
	if err != nil {
		return SubAccount{}, err
	}
	// Unmarshal the body into a payload
	subAccount := SubAccount{}
	err = json.Unmarshal(body, &subAccount)
	if err != nil {
		return SubAccount{}, err
	}
	// Validate the status of the sub account
	if subAccount.Status != "ACTIVE" {
		return SubAccount{}, errors.New("From account is not active")
	}
	// Return the sub account
	return subAccount, nil
}

func (exchangeOrder *ExchangeOrder) ValidateOrder() error {
	// Validate the from sub account
	fromSubAccount, err := ValidateSubAccountStatus(exchangeOrder.From)
	if err != nil {
		return err
	}
	// Validate the from account balance
	if fromSubAccount.Balance < exchangeOrder.Amount {
		return errors.New("Insufficient balance")
	}
	// Validate the fund sub account
	_, err = ValidateSubAccountStatus(exchangeOrder.Fund)
	if err != nil {
		return err
	}

	return nil
}

// BeforeInit is a function that is called before the exchange order is initialized
func (exchangeOrder *ExchangeOrder) BeforeInit() {
	exchangeOrder.ReceivedAt = time.Now()
	exchangeOrder.TransactionId = uuid.New()
}

// InitExchangeRequestHandling is a function that is called after the exchange order is initialized
func (exchangeOrder *ExchangeOrder) InitExchangeRequestHandling() error {
	exchangeOrder.BeforeInit()

	err := exchangeOrder.ValidateOrder()
	if err != nil {
		return err
	}

	exchangeOrder.ValidatedAt = time.Now()
	return nil
}
