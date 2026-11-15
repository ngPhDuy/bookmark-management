package redis

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	RedisHost     string `default:"localhost" envconfig:"REDIS_HOST"`
	RedisPort     string `default:"6379" envconfig:"REDIS_PORT"`
	RedisPassword string `default:"" envconfig:"REDIS_PASSWORD"`
	RedisDB       int    `default:"0" envconfig:"REDIS_DB"`
}

func NewRedisConfig() *RedisConfig {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	if err := godotenv.Load(); err != nil {
		logger.Warn("Load redis config fail", "err", err)
	}

	cfg := &RedisConfig{}

	//stop if process config fail, because we don't want to run without config
	if err := envconfig.Process("", cfg); err != nil {
		logger.Error("Process redis config fail", "err", err)
		os.Exit(1)
	}

	return cfg
}
