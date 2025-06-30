package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"go-auth/internal/jwt"
	"go-auth/internal/models"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenRequest struct {
	UserID string `json:"user_id"`
}

// TokensHandler POST /tokens?user_id=GUID
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
	refreshRaw := make([]byte, 32)
	_, err = rand.Read(refreshRaw)
	if err != nil {
		http.Error(w, "Ошибка генерации refresh токена", http.StatusInternalServerError)
		return
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(refreshRaw)
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
	json.NewEncoder(w).Encode(resp)
}
