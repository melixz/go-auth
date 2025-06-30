package handlers

import (
	"go-auth/internal/jwt"
	"net/http"
	"strings"
)

// LogoutHandler POST /logout (деавторизация)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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
	userID := claims.UserID
	_ = userID
	w.WriteHeader(http.StatusNoContent)
}
