package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AppPort       int    `default:"8080" envconfig:"APP_PORT"`
	ServiceName   string `default:"bookmark-management-default" envconfig:"SERVICE_NAME"`
	InstanceID    string `default:"" envconfig:"INSTANCE_ID"`
	RedisHost     string `default:"localhost" envconfig:"REDIS_HOST"`
	RedisPort     int    `default:"6379" envconfig:"REDIS_PORT"`
	RedisPassword string `default:"" envconfig:"REDIS_PASSWORD"`
}

// NewAppConfig creates a new AppConfig instance by loading environment variables
// from a .env file and processing them into the AppConfig struct.
func NewAppConfig() *AppConfig {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	if err := godotenv.Load(); err != nil {
		logger.Warn("Load config fail", "err", err)
	}

	cfg := &AppConfig{}

	err := envconfig.Process("", cfg)

	if err != nil {
		logger.Warn("Process config fail", "err", err)
	}

	return cfg
}
