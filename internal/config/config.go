package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Config содержит все настройки приложения
type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ExternalAPIURL string
	LogLevel       string
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env файла: %w", err)
	}

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "song_library"),
		ExternalAPIURL: getEnv("EXTERNAL_API_URL", "http://localhost:8081"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
