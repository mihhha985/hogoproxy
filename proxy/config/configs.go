package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DaDataAPIKey    string
	DaDataSecretKey string
	JwtSecret       string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	godotenv.Load()

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		panic("Can not get secret")
	}

	apiKey := os.Getenv("DADATA_API_KEY")
	if apiKey == "" {
		panic("Can not get DADATA_API_KEY")
	}

	dadataSecretKey := os.Getenv("DADATA_SECRET_KEY")
	if dadataSecretKey == "" {
		panic("Can not get DADATA_SECRET_KEY")
	}

	return &Config{
		DaDataAPIKey:    apiKey,
		DaDataSecretKey: dadataSecretKey,
		JwtSecret:       secret,
	}
}
