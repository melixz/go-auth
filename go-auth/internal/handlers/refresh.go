package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"

	"go-auth/internal/db"
	"go-auth/internal/jwt"
	"go-auth/internal/models"
	"go-auth/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshHandler обновляет пару токенов
// @Summary      Обновить пару токенов
// @Description  Принимает действующую пару access+refresh, выдаёт новую пару
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Пара токенов"
// @Success      200  {object}  TokenPairResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tokens/refresh [post]
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	conn, err := db.InitDB(ctx)
	if err != nil {
		http.Error(w, "Ошибка подключения к БД", http.StatusInternalServerError)
		return
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, ctx)

	storedToken, err := db.FindRefreshToken(ctx, conn, userID)
	if err != nil {
		http.Error(w, "Refresh токен не найден или отозван", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(storedToken.BcryptHash), []byte(req.RefreshToken)) != nil {
		http.Error(w, "Неверный refresh токен", http.StatusUnauthorized)
		return
	}
	if storedToken.Used || storedToken.Revoked {
		http.Error(w, "Refresh токен уже использован или отозван", http.StatusUnauthorized)
		return
	}
	if storedToken.UserAgent != userAgent {
		_ = db.RevokeRefreshToken(ctx, conn, storedToken.ID)
		http.Error(w, "User-Agent не совпадает, токен отозван", http.StatusUnauthorized)
		return
	}
	if storedToken.IP != ip {
		utils.SendWebhook(userID, storedToken.IP, ip, userAgent)
	}
	_ = db.MarkRefreshTokenUsed(ctx, conn, storedToken.ID)

	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Не удалось сгенерировать токен доступа", http.StatusInternalServerError)
		return
	}
	refreshToken, err := utils.GenerateRandomBase64(32)
	if err != nil {
		http.Error(w, "Не удалось сгенерировать refresh токен", http.StatusInternalServerError)
		return
	}
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Не удалось хешировать refresh токен", http.StatusInternalServerError)
		return
	}
	newToken := models.RefreshToken{
		UserID:     userID,
		BcryptHash: string(refreshHash),
		UserAgent:  userAgent,
		IP:         ip,
		CreatedAt:  time.Now(),
		Revoked:    false,
		Used:       false,
	}
	if err := db.CreateRefreshToken(ctx, conn, newToken); err != nil {
		http.Error(w, "Не удалось сохранить refresh токен", http.StatusInternalServerError)
		return
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
