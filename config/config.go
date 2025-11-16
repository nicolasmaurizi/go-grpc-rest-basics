package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AppConfig struct {
	DB       DBConfig
	GRPCPort string
	HTTPPort string
}

// helper interno
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() AppConfig {
	dbCfg := DBConfig{
		Host:     getenv("PG_HOST", "localhost"),
		Port:     getenv("PG_PORT", "5432"),
		User:     getenv("PG_USER", "postgres"),
		Password: getenv("PG_PASSWORD", "admin"),
		Name:     getenv("PG_DB", "bloomgrpc"),
		SSLMode:  getenv("PG_SSLMODE", "disable"),
	}

	return AppConfig{
		DB:       dbCfg,
		GRPCPort: getenv("GRPC_PORT", "50051"),
		HTTPPort: getenv("HTTP_PORT", "8080"),
	}
}

// Método útil para armar el connection string
func (d DBConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}
