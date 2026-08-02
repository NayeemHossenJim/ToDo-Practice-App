package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPAddress string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	applicationConfig := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		HTTPAddress: os.Getenv("HTTP_ADDRESS"),
	}

	if applicationConfig.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	if applicationConfig.HTTPAddress == "" {
		applicationConfig.HTTPAddress = ":8080"
	}

	return applicationConfig, nil
}
