package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		data, err := json.Marshal(map[string]string{"token": ""})
		if err != nil {
			slog.Error("signIn: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(data); err != nil {
			slog.Error("signIn: failed to write response", "error", err)
		}
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		data, err := json.Marshal(map[string]string{"error": "Неверный пароль"})
		if err != nil {
			slog.Error("signIn: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(data); err != nil {
			slog.Error("signIn: failed to write response", "error", err)
		}
		return
	}

	if req.Password != pass {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		data, err := json.Marshal(map[string]string{"error": "Неверный пароль"})
		if err != nil {
			slog.Error("signIn: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(data); err != nil {
			slog.Error("signIn: failed to write response", "error", err)
		}
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": hash,
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(pass))
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	data, err := json.Marshal(map[string]string{"token": tokenStr})
	if err != nil {
		slog.Error("signIn: failed to marshal response", "error", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		slog.Error("signIn: failed to write response", "error", err)
	}
}
