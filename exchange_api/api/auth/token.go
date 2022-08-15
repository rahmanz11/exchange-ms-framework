package auth

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// CreateToken creates a new JWT token
func CreateToken(subAccountId uuid.UUID) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["sub_account_id"] = subAccountId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token and return
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// TokenValid checks if the token is valid
func TokenValid(r *http.Request) error {
	// Extract token from request
	tokenString := ExtractToken(r)
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	// Validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}

	return nil
}

// ExtractToken extracts the token from the request
func ExtractToken(r *http.Request) string {
	// Get the token from the header
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	// Get the token from the auth header
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts the token ID from the request
func ExtractTokenID(r *http.Request) (string, error) {
	// Extract token from request
	tokenString := ExtractToken(r)
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	// Validate token
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if err != nil {
			return "", err
		}
		// Return account ID
		return claims["sub_account_id"].(string), nil
	}
	return "", nil
}

// Pretty display the claims
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
