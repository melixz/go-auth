package db

import (
	"context"
	"go-auth/internal/models"

	"github.com/jackc/pgx/v5"
)

func CreateRefreshToken(ctx context.Context, conn *pgx.Conn, token models.RefreshToken) error {
	_, err := conn.Exec(ctx, `INSERT INTO refresh_tokens (user_id, bcrypt_hash, user_agent, ip, created_at, revoked, used) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		token.UserID, token.BcryptHash, token.UserAgent, token.IP, token.CreatedAt, token.Revoked, token.Used)
	return err
}

func FindRefreshToken(ctx context.Context, conn *pgx.Conn, userID string) (*models.RefreshToken, error) {
	row := conn.QueryRow(ctx, `SELECT id, user_id, bcrypt_hash, user_agent, ip, created_at, revoked, used FROM refresh_tokens WHERE user_id=$1 AND revoked=false AND used=false ORDER BY created_at DESC LIMIT 1`, userID)
	var t models.RefreshToken
	err := row.Scan(&t.ID, &t.UserID, &t.BcryptHash, &t.UserAgent, &t.IP, &t.CreatedAt, &t.Revoked, &t.Used)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func RevokeRefreshToken(ctx context.Context, conn *pgx.Conn, id int) error {
	_, err := conn.Exec(ctx, `UPDATE refresh_tokens SET revoked=true WHERE id=$1`, id)
	return err
}

func MarkRefreshTokenUsed(ctx context.Context, conn *pgx.Conn, id int) error {
	_, err := conn.Exec(ctx, `UPDATE refresh_tokens SET used=true WHERE id=$1`, id)
	return err
}
