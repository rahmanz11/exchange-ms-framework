package auth

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// CreateToken creates a new token for a user using sub_account_id
func CreateToken(subAccountId uuid.UUID) (string, error) {
	// Create the token claims
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["sub_account_id"] = subAccountId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token and return
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// TokenValid checks if a token is valid
func TokenValid(r *http.Request) error {
	// Extract token from request
	tokenString := ExtractToken(r)
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	// If there is an error, return it
	if err != nil {
		return err
	}
	// Is token valid?
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Print the token in a pretty way for logging puposes
		Pretty(claims)
	}
	
	return nil
}

// ExtractToken extracts the token from the request
func ExtractToken(r *http.Request) string {
	// Get the token from the request header
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	// If a valid token is found, return it
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	// If there is no token, return an empty string
	return ""
}

// ExtractTokenID returns the token ID from a token
func ExtractTokenID(r *http.Request) (uint32, error) {
	// Get token from request
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	// Extract token ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

// Pretty prints the token claims in a pretty way
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
