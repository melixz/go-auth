package handlers

import (
	"encoding/json"
	"go-auth/internal/jwt"
	"net/http"
	"strings"
)

type MeResponse struct {
	UserID string `json:"user_id"`
}

// MeHandler GET /me (защищённый endpoint)
func MeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Неверный формат заголовка авторизации", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ValidateAccessToken(token)
	if err != nil {
		http.Error(w, "Неверный токен доступа", http.StatusUnauthorized)
		return
	}
	resp := MeResponse{
		UserID: claims.UserID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
