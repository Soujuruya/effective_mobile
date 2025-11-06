package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"effective_mobile/internal/config"
	"effective_mobile/internal/migrator"
	"effective_mobile/internal/repository"
	"effective_mobile/internal/service"
	v1 "effective_mobile/internal/transport/http/v1"
	"effective_mobile/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const migrationsDir = "migrations"

func main() {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env"
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Создаём мигратор
	migr := migrator.MustGetNewMigrator(MigrationsFS, migrationsDir)

	dbURL := cfg.BuildDatabaseURL()

	// Настраиваем логгер
	lg := logger.NewLogger(cfg.Environment)
	defer lg.Sync()

	ctx := context.Background()
	if err := migr.ApplyMigrations(dbURL); err != nil {
		lg.Error(ctx, "failed to apply migrations", zap.Error(err))
	}

	// Подключаемся к базе через pgxpool
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		lg.Error(ctx, "failed to connect to database: %v", zap.Error(err))
		return
	}
	defer db.Close()

	// Инициализация репозитория и сервиса
	repoSubs := repository.NewRepository(db, cfg.Environment)
	subsService := service.NewSubscriptionService(repoSubs, cfg.Environment)

	server := v1.NewServer(cfg.Port, subsService)
	server.RegisterHandlers()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		lg.Info(ctx, "HTTP server listening on port %d", zap.Int("port", cfg.Port))
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			lg.Error(ctx, "server error: %v", zap.Error(err))
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	lg.Info(ctx, "Shutdown signal received, starting graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		lg.Info(ctx, "server shutdown error: %v", zap.Error(err))
	}

	lg.Info(ctx, "Database connection pool closed")
	wg.Wait()
	lg.Info(ctx, "Server stopped gracefully")
}
