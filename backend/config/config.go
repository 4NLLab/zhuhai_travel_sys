package config

import (
	"os"
)

type Config struct {
	ServerPort           string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	TokenTTLHours        string
	CORSAllowedOrigin    string
	PaymentWebhookSecret string
}

func Load() *Config {
	return &Config{
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "127.0.0.1"),
		DBPort:               getEnv("DB_PORT", "3306"),
		DBUser:               getEnv("DB_USER", "root"),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBName:               getEnv("DB_NAME", "zhuhai_travel"),
		JWTSecret:            getEnv("JWT_SECRET", "change-me-in-production"),
		TokenTTLHours:        getEnv("TOKEN_TTL_HOURS", "168"),
		CORSAllowedOrigin:    getEnv("CORS_ALLOWED_ORIGIN", "*"),
		PaymentWebhookSecret: getEnv("PAYMENT_WEBHOOK_SECRET", "change-me-payment-webhook-secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
