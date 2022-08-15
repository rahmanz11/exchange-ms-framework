package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// ExchangeOrder model
type ExchangeOrder struct {
	TransactionId     string    `json:"transaction_id"`
	From              string    `json:"from"`
	Fund              string    `json:"fund"`
	Amount            float64   `json:"amt"`
	Re                string    `json:"re"`
	Status            string    `json:"status"`
	ReceivedAt        time.Time `json:"receivedAt"`
	ValidatedAt       time.Time `json:"validatedAt"`
	FundReceivedAt    time.Time `json:"fund_received_at"`
	FundRequestMadeAt time.Time `json:"fund_request_made_at"`
}

// ValidateSubAccountStatus validates the sub account status
func ValidateSubAccountStatus(accountNumber string) (SubAccount, error) {
	// Get the sub account
	subAccountRaw, err := http.Get(os.Getenv("SUB_ACCOUNT_API") + "/api/v1/" + accountNumber)
	if err != nil {
		return SubAccount{}, err
	}
	// Parse the sub account
	body, err := ioutil.ReadAll(subAccountRaw.Body)
	if err != nil {
		return SubAccount{}, err
	}
	// Unmarshal the sub account
	subAccount := SubAccount{}
	err = json.Unmarshal(body, &subAccount)
	if err != nil {
		return SubAccount{}, err
	}
	// Check the sub account status
	if subAccount.Status != "ACTIVE" {
		return SubAccount{}, errors.New("Account is not active")
	}

	return subAccount, nil
}

// ValidateOrder validates the order
func (exchangeOrder *ExchangeOrder) ValidateOrder() (SubAccount, SubAccount, error) {
	// Validate the from sub account
	fromSubAccount, err := ValidateSubAccountStatus(exchangeOrder.From)
	if err != nil {
		return SubAccount{}, SubAccount{}, err
	}
	// Validate the balance of from sub account
	if fromSubAccount.Balance < exchangeOrder.Amount {
		return SubAccount{}, SubAccount{}, errors.New("Insufficient balance")
	}
	// Validate the fund sub account
	fundSubAccount, err := ValidateSubAccountStatus(exchangeOrder.Fund)
	if err != nil {
		return SubAccount{}, SubAccount{}, err
	}
	// return the sub accounts
	return fromSubAccount, fundSubAccount, nil
}

// BeforeInit initializes the exchange order
func (exchangeOrder *ExchangeOrder) BeforeInit() {
	exchangeOrder.FundReceivedAt = time.Now()
}

// ExchangeOrderRequestHandling handles the exchange order request
func (exchangeOrder *ExchangeOrder) ExchangeRequestHandling() error {
	// Prepare the exchange order request
	exchangeOrder.BeforeInit()
	// Validate the exchange order
	fromSubAccount, fundSubAccount, err := exchangeOrder.ValidateOrder()
	if err != nil {
		return err
	}
	// Set timestamps for the exchange order
	exchangeOrder.ValidatedAt = time.Now()
	exchangeOrder.FundRequestMadeAt = time.Now()

	// Update the balance of the from sub account
	_, err = http.Post(os.Getenv("SUB_ACCOUNT_API")+"/api/v1/update/"+exchangeOrder.From,
		"application/json", bytes.NewBuffer([]byte(`{"balance":`+fmt.Sprintf("%.2f", fromSubAccount.Balance-exchangeOrder.Amount)+`}`)))
	if err != nil {
		return err
	}
	// Update the balance of the fund sub account
	_, err = http.Post(os.Getenv("SUB_ACCOUNT_API")+"/api/v1/update/"+exchangeOrder.Fund,
		"application/json", bytes.NewBuffer([]byte(`{"balance":`+fmt.Sprintf("%.2f", fundSubAccount.Balance+exchangeOrder.Amount)+`}`)))
	// Set the status of the exchange order
	exchangeOrder.Status = "PENDING"
	return nil
}
