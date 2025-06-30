package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"go-auth/internal/jwt"
	"go-auth/internal/models"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshHandler POST /tokens/refresh
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if req.AccessToken == "" || req.RefreshToken == "" {
		http.Error(w, "access_token и refresh_token обязательны", http.StatusBadRequest)
		return
	}
	claims, err := jwt.ValidateAccessToken(req.AccessToken)
	if err != nil {
		http.Error(w, "Неверный токен доступа", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID
	userAgent := r.UserAgent()
	ip := r.RemoteAddr
	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Не удалось сгенерировать токен доступа", http.StatusInternalServerError)
		return
	}
	refreshRaw := make([]byte, 32)
	_, err = rand.Read(refreshRaw)
	if err != nil {
		http.Error(w, "Не удалось сгенерировать refresh токен", http.StatusInternalServerError)
		return
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(refreshRaw)
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Не удалось хешировать refresh токен", http.StatusInternalServerError)
		return
	}
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

func sendWebhook(userID, oldIP, newIP, userAgent string) {
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return
	}
	payload := map[string]interface{}{
		"user_id":    userID,
		"old_ip":     oldIP,
		"new_ip":     newIP,
		"user_agent": userAgent,
		"timestamp":  time.Now().Unix(),
	}
	jsonData, _ := json.Marshal(payload)
	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
}
