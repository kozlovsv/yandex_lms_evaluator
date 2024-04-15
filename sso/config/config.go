package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	JWT  JWTConfig
	DB   DBConfig
	HTTP HttpConfig
}

type JWTConfig struct {
	Secret   string
	TokenTTL time.Duration
}

type HttpConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
	Params   string
}

func NewConfig() *Config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found. ")
	}

	return &Config{
		Env: getEnv("env", "local"),
		HTTP: HttpConfig{
			Port: getEnv("SSO_HTTP_PORT", ""),
			Host: getEnv("SSO_HTTP_HOST", ""),
		},
		JWT: JWTConfig{
			Secret:   getEnv("JWT_SECRET", "sekret_key"),
			TokenTTL: getEnvAsDuration("JWT_TOKEN_TTL", time.Hour*24),
		},
		DB: DBConfig{
			User:     getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", "3306"),
			DbName:   getEnv("DB_NAME", ""),
			Params:   getEnv("DB_CONNACTION_PARAMS", ""),
		},
	}
}

func (d *DBConfig) GetDSNString() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Password, d.Host, d.Port, d.DbName)
	if d.Params != "" {
		dsn += "?" + d.Params
	}
	return dsn
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
/*func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}*/

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(name, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}
