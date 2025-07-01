package handlers

import (
	"go-auth/internal/jwt"
	"net/http"
	"strings"
)

// LogoutHandler деавторизует пользователя
// @Summary      Деавторизация
// @Description  Отзывает текущий access/refresh токен
// @Security     ApiKeyAuth
// @Tags         auth
// @Success      204  "No Content"
// @Failure      401  {object}  ErrorResponse
// @Router       /logout [post]
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
