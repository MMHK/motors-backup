package config

import (
	"github.com/joho/godotenv"
	"motors-backup/internal/log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Config holds the configuration for database connection
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	port, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "3306"))
	if err != nil {
		port = 3306
	}

	return &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     port,
		DBUser:     getEnvOrDefault("DB_USER", "root"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:     os.Getenv("DB_NAME"), // DB_NAME is required, no default value
	}
}

func GetLocalPath(file string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), file)
}

func LoadTestConfig() *Config {
	err := godotenv.Load(GetLocalPath("../../.env"))
	if err != nil {
		log.Logger.Errorf("Error loading environment: %s", err)
	}

	return &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     3306,
		DBUser:     getEnvOrDefault("DB_USER", "root"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:     os.Getenv("DB_NAME"), // DB_NAME is required, no default value
	}
}
