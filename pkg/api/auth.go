package api

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenStr string
			if cookie, err := r.Cookie("token"); err == nil {
				tokenStr = cookie.Value
			}
			if !validateToken(tokenStr, pass) {
				slog.Warn("unauthorized request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func validateToken(tokenStr, pass string) bool {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(pass), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hash, ok := claims["password_hash"].(string)
	if !ok {
		return false
	}
	expected := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
	return hash == expected
}
