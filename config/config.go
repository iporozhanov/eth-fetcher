package config

import (
	"os"
	"time"
)

type Config struct {
	APIPort         string
	ETHNodeURL      string
	DBConnectionURL string
	JWT             JWT
}

func LoadConfig() *Config {
	jwt := JWT{}
	jwt.Default()

	if os.Getenv("JWT_SECRET") != "" {
		jwt.Secret = os.Getenv("JWT_SECRET")
	}

	if os.Getenv("JWT_DURATION") != "" {
		duration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
		if err != nil {
			jwt.Duration = duration
		}
	}

	return &Config{
		APIPort:         os.Getenv("API_PORT"),
		ETHNodeURL:      os.Getenv("ETH_NODE_URL"),
		DBConnectionURL: os.Getenv("DB_CONNECTION_URL"),
		JWT:             jwt,
	}
}

type JWT struct {
	Secret   string
	Duration time.Duration
}

func (j *JWT) Default() {
	j.Secret = "secret"
	j.Duration = time.Hour * 24
}
