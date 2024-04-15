package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	DB       DBConfig
	GRPCPort string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
}

func NewConfig() *Config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found. ")
	}

	return &Config{
		Env: getEnv("env", "local"),
		DB: DBConfig{
			User:     getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", "3306"),
			DbName:   getEnv("DB_NAME", ""),
		},
		GRPCPort: getEnv("SERVER_GRPC_PORT", "50002"),
	}
}

func (d *DBConfig) GetDSNString() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", d.User, d.Password, d.Host, d.Port, d.DbName)
	return dsn
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
