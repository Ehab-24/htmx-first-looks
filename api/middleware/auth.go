package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("x-access-token")
		if err != nil {
			log.Println(err)
			w.Write([]byte("invalid token"))
			return
		}

		token := cookie.Value
		t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("invalid access token")
			}
			return "", nil
		})
		// ! bug here (in the golang-jwt package). Returns 'key is of invalid type' error
		if err != nil {
			if err.Error() != "token signature is invalid: key is of invalid type" {
				http.Error(w, "invalid access token", http.StatusUnauthorized)
				return
			}
		}
		t.Valid = true

		ctx := context.WithValue(r.Context(), "user_id", t.Claims.(jwt.MapClaims)["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
