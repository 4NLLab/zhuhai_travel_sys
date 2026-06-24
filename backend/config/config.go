package config

import (
	"bufio"
	"os"
	"strings"
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
	IslandCruiseBaseURL  string
	IslandDistributor    string
	IslandAccessToken    string
}

func Load() *Config {
	loadEnvFile(".env")
	loadEnvFile("backend/.env")

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
		IslandCruiseBaseURL:  strings.TrimRight(getEnv("ISLAND_CRUISE_BASE_URL", "http://121.46.23.81:8035"), "/"),
		IslandDistributor:    getEnv("ISLAND_CRUISE_DISTRIBUTOR_CODE", ""),
		IslandAccessToken:    getEnv("ISLAND_CRUISE_ACCESS_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		_ = os.Setenv(key, value)
	}
}
