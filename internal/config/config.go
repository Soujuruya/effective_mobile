package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`

	Port    int           `env:"PORT" env-default:"8080"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`

	BaseURL string `env:"BASE_URL" env-default:"http://localhost:8080"`

	// PostgreSQL
	DBUser     string `env:"DB_USER" env-default:"appuser"`
	DBPassword string `env:"DB_PASSWORD" env-default:"123"`
	DBHost     string `env:"DB_HOST" env-default:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT" env-default:"5432"`
	DBName     string `env:"DB_NAME" env-default:"effective_mobile"`
	DBSSLMode  string `env:"DB_SSLMODE" env-default:"disable"`

	// Полный URL для подключения (можно использовать вместо отдельных полей)
	DatabaseURL string `env:"DATABASE_URL"`
}

// BuildDatabaseURL возвращает полный URL подключения к PostgreSQL.
// Если DATABASE_URL задан в .env, используется он, иначе собирается из отдельных полей.
func (c *Config) BuildDatabaseURL() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
