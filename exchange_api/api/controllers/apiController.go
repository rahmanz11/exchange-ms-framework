package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Nabeegh-Ahmed/exchange_api/api/models"
	"github.com/Nabeegh-Ahmed/exchange_api/api/responses"
)

// CreateSubAccount creates a sub account
func (server *Server) CreateSubAccount(res http.ResponseWriter, req *http.Request) {
	// Get the body of our POST request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Post(os.Getenv("SUB_ACCOUNT_API")+"/api/v1/register",
		"application/json", bytes.NewBuffer(body))
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	// Decode the response into the map
	err = json.NewDecoder(request.Body).Decode(&httpResponse)

	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}

	// If there is no error, return the response
	type apiResponse struct {
		SubAccountUUID string `json:"sub_account_uuid"`
	}
	responses.JSON(res, http.StatusCreated, apiResponse{SubAccountUUID: httpResponse["sub_account_uuid"].(string)})
}

// Login is a function to login a sub account
func (server *Server) Login(res http.ResponseWriter, req *http.Request) {
	// Get the body of our POST request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Post(os.Getenv("SUB_ACCOUNT_API")+"/api/v1/login",
		"application/json", bytes.NewBuffer(body))

	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// If there is no error, return the response
	type apiResponse struct {
		Token string `json:"token"`
	}
	responses.JSON(res, http.StatusOK, apiResponse{Token: httpResponse["token"].(string)})
}

// LinkAccount is a function to link accounts to a sub account
func (server *Server) LinkAccount(res http.ResponseWriter, req *http.Request) {
	// Get the body of our POST request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Post(os.Getenv("SUB_ACCOUNT_API")+"/api/v1/link-account", "application/json", bytes.NewBuffer(body))
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// If there is no error, return the response
	responses.JSON(res, request.StatusCode, map[string]string{"message": "Account linked successfully"})
}

// Exchange handles exhange requests
func (server *Server) Exchange(res http.ResponseWriter, req *http.Request) {
	// Get the body of our POST request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	exchangeOrder := models.ExchangeOrder{}
	err = json.Unmarshal(body, &exchangeOrder)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the validate api
	request, err := http.Post(os.Getenv("VALIDATE_API")+"/api/v1/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	// Decode the response into the map
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// If there is no error, return the response
	responses.JSON(res, request.StatusCode, httpResponse)
}

// WireIn
func (server *Server) WireIn(res http.ResponseWriter, req *http.Request) {

}

// WireOut
func (server *Server) WireOut(res http.ResponseWriter, req *http.Request) {

}

// GetStatus returns the status of the sub account
func (server *Server) GetStatus(res http.ResponseWriter, req *http.Request) {
	// Get the body of our GET request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	apiRequest := models.AccountNumberPayload{}
	err = json.Unmarshal(body, &apiRequest)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Get(os.Getenv("SUB_ACCOUNT_API") + "/api/v1/" + apiRequest.AccountNumber)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// verify the account number is the same as the one in the token
	if httpResponse["sub_account_id"].(string) != req.Context().Value("SubAccountId") {
		responses.ERROR(res, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// If there is no error, return the response
	type apiResponse struct {
		Status string `json:"status"`
	}
	responses.JSON(res, http.StatusOK, apiResponse{Status: httpResponse["status"].(string)})
}

// EnableAccount is a function to set ACTIVE status of an account
func (server *Server) EnableAccount(res http.ResponseWriter, req *http.Request) {
	EnableDisableHelper(res, req, "ACTIVE")
}

// DisableAccount is a function to set HOLD status of an account
func (server *Server) DisableAccount(res http.ResponseWriter, req *http.Request) {
	EnableDisableHelper(res, req, "HOLD")
}

// EnableDisableHelper is a helper function to enable or disable an account
func EnableDisableHelper(res http.ResponseWriter, req *http.Request, status string) {
	// Get the body of our POST request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	apiRequest := models.AccountNumberPayload{}
	err = json.Unmarshal(body, &apiRequest)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Get(os.Getenv("SUB_ACCOUNT_API") + "/api/v1/" + apiRequest.AccountNumber)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// verify the account number is the same as the one in the token
	if httpResponse["sub_account_id"].(string) != req.Context().Value("SubAccountId") {
		responses.ERROR(res, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Marshal the body into a payload
	updatePayload, err := json.Marshal(map[string]string{"status": status})
	// Send an update request to the sub account api
	request, err = http.Post(os.Getenv("SUB_ACCOUNT_API")+"/api/v1/update/"+apiRequest.AccountNumber, "application/json", bytes.NewBuffer(updatePayload))
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	responses.JSON(res, http.StatusOK, map[string]string{"status": status})
}

func (server *Server) GetTransaction(res http.ResponseWriter, req *http.Request) {
	// Get the body of our GET request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	apiRequest := models.TransactionIdPayload{}
	err = json.Unmarshal(body, &apiRequest)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Get(os.Getenv("TRANSACTIONS_API") + "/api/v1/" + apiRequest.TransactionId)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api to get the account details of the account that originated the request
	subAccountId := req.Context().Value("SubAccountId").(string)
	subAccountRequest, err := http.Get(os.Getenv("SUB_ACCOUNT_API") + "/api/v1/id/" + subAccountId)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	// Decode the response into a map
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// Define a new map to hold the response
	var subAccount map[string]interface{}
	err = json.NewDecoder(subAccountRequest.Body).Decode(&subAccount)
	// Check if the token account id is the same as the one in the from field or fund field in the transaction
	if subAccount["account_number"] != httpResponse["from"] && subAccount["account_number"] != httpResponse["fund"] {
		responses.ERROR(res, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// If there is no error, return the response
	responses.JSON(res, http.StatusOK, map[string]any{"transaction": httpResponse, "balance": subAccount["balance"]})
}

// TransactionsRegister returns a list of transactions for a sub account
func (server *Server) TransactionsRegister(res http.ResponseWriter, req *http.Request) {
	// Get the body of our GET request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	apiRequest := models.AccountNumberPayload{}
	err = json.Unmarshal(body, &apiRequest)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the Transactions api
	request, err := http.Get(os.Getenv("TRANSACTIONS_API") + "/api/v1/transactions/" + apiRequest.AccountNumber)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// If there is no error, return the response
	responses.JSON(res, request.StatusCode, httpResponse)
}

// GetTransactionConfirmation gets status of a transaction
func (server *Server) GetTransactionConfirmation(res http.ResponseWriter, req *http.Request) {
	// Get the body of our GET request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	apiRequest := models.AccountNumberPayload{}
	err = json.Unmarshal(body, &apiRequest)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the Transactions api
	request, err := http.Get(os.Getenv("TRANSACTIONS_API") + "/api/v1/" + apiRequest.AccountNumber)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// If there is no error, return the response
	responses.JSON(res, http.StatusOK, map[string]string{"status": httpResponse["status"].(string)})
}

// GetBalance returns the balance of a sub account
func (server *Server) GetBalance(res http.ResponseWriter, req *http.Request) {
	// Get the body of our GET request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	apiRequest := models.AccountNumberPayload{}
	err = json.Unmarshal(body, &apiRequest)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Send the request to the sub account api
	request, err := http.Get(os.Getenv("SUB_ACCOUNT_API") + "/api/v1/" + apiRequest.AccountNumber)
	// If there is an error, return error
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Define a new map to hold the response
	var httpResponse map[string]interface{}
	err = json.NewDecoder(request.Body).Decode(&httpResponse)
	// If there is an error, return error
	if httpResponse["error"] != nil {
		responses.ERROR(res, http.StatusInternalServerError, errors.New(httpResponse["error"].(string)))
		return
	}
	// Check if the account that made the request is the same as the one getting balance for
	if httpResponse["sub_account_id"].(string) != req.Context().Value("SubAccountId") {
		responses.ERROR(res, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// If there is no error, return the response
	responses.JSON(res, http.StatusOK, map[string]float64{"balance": httpResponse["balance"].(float64)})

}
