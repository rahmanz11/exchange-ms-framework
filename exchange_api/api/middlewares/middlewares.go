package middlewares

import (
	"context"
	"errors"
	"github.com/Nabeegh-Ahmed/exchange_api/api/auth"
	"github.com/Nabeegh-Ahmed/exchange_api/api/responses"
	"net/http"
)

// SetMiddlewareJSON is a middleware that parses the request body as JSON
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication is a middleware that checks if the request has a valid token
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		id, err := auth.ExtractTokenID(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		authCtx := context.WithValue(r.Context(), "SubAccountId", id)
		next.ServeHTTP(w, r.WithContext(authCtx))
	}
}
