package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-auth/internal/handlers"
	"go-auth/internal/middleware"

	_ "go-auth/docs"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Auth API
// @version 1.0
// @description Сервис аутентификации (Go, PostgreSQL, Docker)
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env не найден, используются переменные окружения")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не задан")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	log.Println("Подключение к БД успешно")

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		if _, err := os.Stat("/app/migrations"); err == nil {
			migrationsPath = "file:///app/migrations"
		} else {
			migrationsPath = "file://./migrations"
		}
	}
	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatalf("ошибка создания мигратора: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("ошибка применения миграций: %v", err)
	}
	log.Println("Миграции применены (или не требуются)")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "ok")
		if err != nil {
			return
		}
	})
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /tokens", func(w http.ResponseWriter, r *http.Request) {
		handlers.TokensHandler(w, r)
	})
	mux.HandleFunc("POST /tokens/refresh", func(w http.ResponseWriter, r *http.Request) {
		handlers.RefreshHandler(w, r)
	})
	mux.Handle("GET /me", middleware.AuthMiddleware(http.HandlerFunc(handlers.MeHandler)))
	mux.Handle("POST /logout", middleware.AuthMiddleware(http.HandlerFunc(handlers.LogoutHandler)))

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", middleware.CORS(mux)); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
