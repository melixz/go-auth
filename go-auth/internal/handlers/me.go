package handlers

import (
	"encoding/json"
	"go-auth/internal/middleware"
	"net/http"
)

// MeResponse содержит GUID пользователя
// swagger:model MeResponse
type MeResponse struct {
	UserID string `json:"user_id" example:"b3e1c2d4-5678-4a9b-8c2d-1234567890ab"`
}

// MeHandler возвращает GUID текущего пользователя
// @Summary      Получить GUID текущего пользователя
// @Description  Защищённый эндпоинт. Требуется Bearer access-токен.
// @Security     ApiKeyAuth
// @Tags         users
// @Produce      json
// @Success      200  {object}  MeResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /me [get]
func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}
	resp := MeResponse{UserID: claims.UserID}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
