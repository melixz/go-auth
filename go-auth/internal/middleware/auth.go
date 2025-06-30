package middleware

import (
	"context"
	"go-auth/internal/jwt"
	"go-auth/internal/models"
	"net/http"
	"strings"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ValidateAccessToken(token)
		if err != nil {
			http.Error(w, "Неверный токен доступа", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(r *http.Request) *models.AccessTokenClaims {
	val := r.Context().Value(ClaimsKey)
	if claims, ok := val.(*models.AccessTokenClaims); ok {
		return claims
	}
	return nil
}
