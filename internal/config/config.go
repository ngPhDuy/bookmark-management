package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AppPort     int    `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-management-default" envconfig:"SERVICE_NAME"`
	InstanceId  string `default:"" envconfig:"INSTANCE_ID"`
}

func NewAppConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &AppConfig{}

	err := envconfig.Process("", cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
