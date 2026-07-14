package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	EnableSwagger bool
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		EnableSwagger: getEnv("SWAGGER_ENABLED", "false") == "true",
	}

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		u, err := url.Parse(dsn)
		if err == nil {
			cfg.DBHost = u.Hostname()
			cfg.DBPort = u.Port()
			if cfg.DBPort == "" {
				cfg.DBPort = "5432"
			}
			cfg.DBUser = u.User.Username()
			cfg.DBPassword, _ = u.User.Password()
			cfg.DBName = strings.TrimPrefix(u.Path, "/")
			if sslmode := u.Query().Get("sslmode"); sslmode != "" {
				cfg.DBSSLMode = sslmode
			} else {
				cfg.DBSSLMode = "require"
			}
		}
	} else {
		cfg.DBHost = getEnv("DATABASE_HOST", "localhost")
		cfg.DBPort = getEnv("DATABASE_PORT", "5432")
		cfg.DBUser = getEnv("DATABASE_USER", "postgres")
		cfg.DBPassword = getEnv("DATABASE_PASSWORD", "postgres")
		cfg.DBName = getEnv("DATABASE_NAME", "menu_tree")
		cfg.DBSSLMode = getEnv("DATABASE_SSLMODE", "disable")
	}

	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
