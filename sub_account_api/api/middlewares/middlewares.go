package middlewares

import (
	"errors"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/auth"
	"github.com/Nabeegh-Ahmed/sub_account_api/api/responses"
	"net/http"
)

// SetMiddlewareJSON makes sure that the request is in JSON format
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication makes sure that the request has a valid token
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
