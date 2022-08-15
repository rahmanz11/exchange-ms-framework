package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Nabeegh-Ahmed/validate/api/models"
	"github.com/Nabeegh-Ahmed/validate/api/responses"
	"github.com/google/uuid"
)

// SubmitExchangeOrder handles the request on /api/v1/
func (server *Server) SubmitExchangeOrder(res http.ResponseWriter, req *http.Request) {
	// Read the body of the request
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
	// Validate the payload and initiate the exchange order
	err = exchangeOrder.InitExchangeRequestHandling()
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	// Marshal the payload into a json string
	exchangeOrderJson, err := json.Marshal(exchangeOrder)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println("before %s", time.Now())
	// Push the exchange order to the kafka
	err = server.KafkaPush(req.Context(), []byte(exchangeOrder.TransactionId.String()), exchangeOrderJson)
	fmt.Println("after %s", time.Now())
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}

	// Return the transaction id in the response
	responses.JSON(res, http.StatusCreated, struct {
		TransactionID uuid.UUID `json:"transaction_id"`
	}{TransactionID: exchangeOrder.TransactionId})
}
