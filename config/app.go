package config

import (
	"os"
)

type AppConfig struct {
	AppPort string
}

func NewAppConfig(isUsingDotEnv bool) AppConfig {

	return AppConfig{
		AppPort: os.Getenv("APP_PORT"),
	}
}
