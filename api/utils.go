package api

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CheckAPIError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println("~ api error:", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return true
	}
	return false
}

func CheckAPIErrorWithStatus(w http.ResponseWriter, err error, message string, code int) bool {
	if err != nil {
		log.Println("~ api error:", err)
		http.Error(w, message, code)
		return true
	}
	return false
}

func generateJWT(secret string, key string, value interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(48 * time.Hour)
	claims["authorized"] = true
	claims[key] = value

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
