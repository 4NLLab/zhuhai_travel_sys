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
	JZGBaseURL           string
	JZGDistributorCode   string
	JZGAccessToken       string
}

func Load() *Config {
	return &Config{
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "127.0.0.1"),
		DBPort:               getEnv("DB_PORT", "3306"),
		DBUser:               getEnv("DB_USER", "root"),
		DBPassword:           getEnv("DB_PASSWORD", "wuyuanjian0"),
		DBName:               getEnv("DB_NAME", "zhuhai_travel"),
		JWTSecret:            getEnv("JWT_SECRET", "zhuhai-travel-secret-key-2026"),
		TokenTTLHours:        getEnv("TOKEN_TTL_HOURS", "168"),
		CORSAllowedOrigin:    getEnv("CORS_ALLOWED_ORIGIN", "*"),
		PaymentWebhookSecret: getEnv("PAYMENT_WEBHOOK_SECRET", "change-me-payment-webhook-secret"),
		JZGBaseURL:           getEnv("JZG_BASE_URL", "http://121.46.23.81:8035"),
		JZGDistributorCode:   getEnv("JZG_DISTRIBUTOR_CODE", ""),
		JZGAccessToken:       getEnv("JZG_ACCESS_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
