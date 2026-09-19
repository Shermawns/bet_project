package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"

	"bet-api/controller"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

func Load() Config {
	return Config{
		DatabaseURL: env(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5432/bet_api?sslmode=disable",
		),
		HTTPAddr: env("HTTP_ADDR", ":8000"),
	}
}

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewDB(lifecycle fx.Lifecycle, config Config) (*sql.DB, error) {
	connection, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return connection.PingContext(ctx)
		},
		OnStop: func(context.Context) error {
			return connection.Close()
		},
	})

	return connection, nil
}

func NewServer(lifecycle fx.Lifecycle, config Config, controller *controller.Controller) *http.Server {
	server := &http.Server{
		Addr:              config.HTTPAddr,
		Handler:           controller.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Println("servidor HTTP:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}
