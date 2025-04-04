package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        int    `envconfig:"PORT" required:"true" validate:"range(1,65535)"`
	Env         string `envconfig:"ENV" required:"true" validate:"oneof=development staging production"`
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true" validate:"startswith=postgresql://"`
	JWTSecret   string `envconfig:"JWT_SECRET" required:"true"`
}

func LoadConfig() (*Config, error) {
	envFile := ".env"
	if customEnvFile := os.Getenv("ENV_FILE"); customEnvFile != "" {
		envFile = customEnvFile
	}

	err := godotenv.Load(envFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
