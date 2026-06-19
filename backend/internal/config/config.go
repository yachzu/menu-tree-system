package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DATABASE_HOST", "localhost"),
		DBPort:     getEnv("DATABASE_PORT", "5432"),
		DBUser:     getEnv("DATABASE_USER", "postgres"),
		DBPassword: getEnv("DATABASE_PASSWORD", "postgres"),
		DBName:     getEnv("DATABASE_NAME", "menu-tree"),
	}
	return cfg
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		"port=" + c.DBPort +
		"user" + c.DBUser +
		"password" + c.DBPassword +
		"dbname" + c.DBName +
		" sslmode=disable"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
