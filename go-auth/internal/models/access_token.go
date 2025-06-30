package models

type AccessTokenClaims struct {
	UserID   string `json:"user_id"`
	IssuedAt int64  `json:"iat"`
	Expiry   int64  `json:"exp"`
}
