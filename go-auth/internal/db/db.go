package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

func InitDB(ctx context.Context) (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_URL")
	return pgx.Connect(ctx, dbURL)
}

func CreateRefreshToken(ctx context.Context, conn *pgx.Conn, token string) error {
	return nil
}

func RevokeRefreshToken(ctx context.Context, conn *pgx.Conn, tokenID int) error {
	return nil
}
