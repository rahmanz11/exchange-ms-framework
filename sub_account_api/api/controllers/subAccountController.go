package controllers

import (
	"encoding/json"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/auth"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/models"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/payloads"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/responses"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/utils/formaterror"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

// Login handles the request on /api/v1/login
/*
req.body {
	@param: subAccountNumber,
	@param: credential
}
*/
func (server *Server) Login(res http.ResponseWriter, req *http.Request) {
	// Read the body of the request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	subAccount := models.SubAccount{}
	err = json.Unmarshal(body, &subAccount)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}

	// Validate the payload
	subAccount.Prepare()
	err = subAccount.Validate("login")
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}

	// SingIn the user using helper function
	token, err := server.SignIn(subAccount.AccountNumber, subAccount.Credential)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(res, http.StatusUnprocessableEntity, formattedError)
		return
	}
	// Return the token in the response
	type loginResponse struct {
		Token string `json:"token"`
	}
	responses.JSON(res, http.StatusOK, loginResponse{Token: token})
}

// SingIn is helper function for login
func (server *Server) SignIn(accountNumber, credential string) (string, error) {
	var err error
	// Check if the account number is valid
	subAccount := models.SubAccount{}
	err = server.db.Debug().Model(models.SubAccount{}).Where("account_number = ?", accountNumber).Take(&subAccount).Error
	if err != nil {
		return "", err
	}
	// Check if the credential is valid
	err = models.VerifyCredential(subAccount.Credential, credential)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	// Create the token and return it
	return auth.CreateToken(subAccount.SubAccountId)
}

// Register handles the request on /api/v1/register
/*
req.body {
	@param: subAccountNumber,
	@param: credential,
	@param: status (optional),
	@param: balance (optional)
}
*/
func (server *Server) Register(res http.ResponseWriter, req *http.Request) {
	// Read the body of the request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
	}
	// Unmarshal the body into a payload
	subAccount := models.SubAccount{}
	err = json.Unmarshal(body, &subAccount)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Create the subAccount
	createdSubAccount, err := subAccount.CreateSubAccount(server.db)
	// If the subAccount is not created, return an error
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(res, http.StatusInternalServerError, formattedError)
		return
	}
	// Return the subAccount in the response
	type response struct {
		SubAccountUUID string `json:"sub_account_uuid"`
	}
	responses.JSON(res, http.StatusCreated, response{SubAccountUUID: createdSubAccount.SubAccountId.String()})
}

// FindSubAccount handles the request on /api/v1/{account-number}
func (server *Server) FindSubAccount(res http.ResponseWriter, req *http.Request) {
	// Get the account number from the request
	params := mux.Vars(req)
	accountNumber := params["account-number"]

	// Find the subAccount
	subAccount := models.SubAccount{}
	foundSubAccount, err := subAccount.FindSubAccountByAccountNumber(server.db, accountNumber)

	// If the subAccount is not found, return an error
	if err != nil {
		responses.ERROR(res, http.StatusNotFound, err)
		return
	}

	// Return the subAccount in the response
	responses.JSON(res, http.StatusOK, foundSubAccount)
}

// UpdateSubAccount handles the request on /api/v1/id/{account-id}
func (server *Server) FindSubAccountById(res http.ResponseWriter, req *http.Request) {
	// Get the subAccount id from the request
	params := mux.Vars(req)
	accountId := params["account-id"]
	// Find the subAccount
	subAccount := models.SubAccount{}
	foundSubAccount, err := subAccount.FindSubAccountByAccountId(server.db, accountId)
	// If the subAccount is not found, return an error
	if err != nil {
		responses.ERROR(res, http.StatusNotFound, err)
		return
	}
	// Return the subAccount in the response
	responses.JSON(res, http.StatusOK, foundSubAccount)
}

// CreateLinkedAccount handles the request on /api/v1/link-account
func (server *Server) CreateLinkedAccount(res http.ResponseWriter, req *http.Request) {
	// Read the body of the request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
	}
	// Unmarshal the body into a payload
	payload := payloads.LinkAccountPayload{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Find the sub account to make sure the account number is valid
	subAccount := models.SubAccount{}
	foundSubAccount, err := subAccount.FindSubAccountByAccountNumber(server.db, payload.SubAccountNumber)
	if err != nil {
		responses.ERROR(res, http.StatusNotFound, err)
		return
	}
	// Create the linked account
	_, err = foundSubAccount.LinkAccount(server.db, payload.AccountNumber)
	if err != nil {
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}
	// Return the subAccount in the response
	responses.JSON(res, http.StatusCreated, "")
}

// UpdateSubAccount handles the request on /api/v1/update/{account-number}
func (server *Server) UpdateSubAccount(res http.ResponseWriter, req *http.Request) {
	// Get the account number from the request
	params := mux.Vars(req)
	accountNumber := params["account-number"]
	// Find the subAccount
	subAccount := models.SubAccount{}
	foundSubAccount, err := subAccount.FindSubAccountByAccountNumber(server.db, accountNumber)
	// If the subAccount is not found, return an error
	if err != nil {
		responses.ERROR(res, http.StatusNotFound, err)
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Unmarshal the body into a payload
	subAccountUpdate := models.SubAccount{}
	err = json.Unmarshal(body, &subAccountUpdate)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Validate the payload and update the subAccount
	if subAccountUpdate.Status == "" {
		subAccountUpdate.Status = foundSubAccount.Status
	}
	if subAccountUpdate.Balance == 0 {
		subAccountUpdate.Balance = foundSubAccount.Balance
	}
	foundSubAccount.Balance = subAccountUpdate.Balance
	foundSubAccount.Status = subAccountUpdate.Status
	// Update the subAccount
	updatedSubAccount, err := foundSubAccount.UpdateSubAccount(server.db)
	if err != nil {
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}
	// Return the updated subAccount in the response
	responses.JSON(res, http.StatusOK, updatedSubAccount)
}
