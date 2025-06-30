package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitDB(ctx context.Context) (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_URL")
	return pgx.Connect(ctx, dbURL)
}
