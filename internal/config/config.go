package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppName     string
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("[Config-1] No .env file found: %v", err)
	}

	return &Config{
		AppName:     viper.GetString("APP_NAME"),
		Port:        viper.GetString("APP_PORT"),
		DatabaseURL: viper.GetString("DB_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		Environment: viper.GetString("APP_ENV"),
	}
}