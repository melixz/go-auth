package handlers

import (
	"encoding/json"
	"go-auth/internal/jwt"
	"go-auth/internal/models"
	"go-auth/internal/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TokenPairResponse содержит пару access и refresh токенов
// swagger:model TokenPairResponse
type TokenPairResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"dGhpc2lzYXJlZnJlc2h0b2tlbg"`
}

type TokenRequest struct {
	UserID string `json:"user_id"`
}

// ErrorResponse описывает ошибку
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"invalid token"`
}

// TokensHandler выдаёт пару токенов
// @Summary      Получить пару токенов
// @Description  Выдаёт access и refresh токены для пользователя с указанным GUID
// @Tags         tokens
// @Param        user_id   query     string  true  "GUID пользователя"
// @Success      200  {object}  TokenPairResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tokens [post]
func TokensHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации access токена", http.StatusInternalServerError)
		return
	}
	refreshToken, err := utils.GenerateRandomBase64(32)
	if err != nil {
		http.Error(w, "Ошибка генерации refresh токена", http.StatusInternalServerError)
		return
	}
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка хеширования refresh токена", http.StatusInternalServerError)
		return
	}
	userAgent := r.UserAgent()
	ip := r.RemoteAddr
	_ = models.RefreshToken{
		UserID:     userID,
		BcryptHash: string(refreshHash),
		UserAgent:  userAgent,
		IP:         ip,
		CreatedAt:  time.Now(),
		Revoked:    false,
		Used:       false,
	}
	resp := TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
