package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv               string
	HTTPPort             int
	DBHost               string
	DBPort               int
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	JWTSecret            string
	JWTAccessTTL         time.Duration
	JWTRefreshTTL        time.Duration
	EncryptionKey        string
	AllowedProviderHosts []string
	RateLimitPerMin      int
	LogLevel             string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		HTTPPort:             getEnvInt("HTTP_PORT", 8080),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnvInt("DB_PORT", 5432),
		DBUser:               getEnv("DB_USER", "sage"),
		DBPassword:           getEnv("DB_PASSWORD", "sage"),
		DBName:               getEnv("DB_NAME", "sage"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTAccessTTL:         time.Duration(getEnvInt("JWT_ACCESS_TTL_MIN", 15)) * time.Minute,
		JWTRefreshTTL:        time.Duration(getEnvInt("JWT_REFRESH_TTL_HOUR", 24*7)) * time.Hour,
		EncryptionKey:        getEnv("ENCRYPTION_KEY", ""),
		AllowedProviderHosts: splitCSV(getEnv("ALLOWED_PROVIDER_HOSTS", "")),
		RateLimitPerMin:      getEnvInt("RATE_LIMIT_PER_MIN", 60),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	if len(c.EncryptionKey) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes")
	}
	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(k string, def int) int {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
