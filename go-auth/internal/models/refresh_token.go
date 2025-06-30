package models

import (
	"time"
)

type RefreshToken struct {
	ID         int       `json:"id"`
	UserID     string    `json:"user_id"`
	BcryptHash string    `json:"-"`
	UserAgent  string    `json:"user_agent"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
	Revoked    bool      `json:"revoked"`
	Used       bool      `json:"used"`
}
