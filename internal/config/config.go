package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	keyHTTPHost = "HTTP_HOST"
	keyHTTPPort = "HTTP_PORT"

	keyDBHost = "DB_HOST"
	keyDBPort = "DB_PORT"
	keyDBUser = "DB_USER"
	keyDBPass = "DB_PASS"
	keyDBName = "DB_NAME"

	LogFatalMissingValue = "%s is missing"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func (cfg *DBConfig) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
}

type HTTPConfig struct {
	Host string
	Port string
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf(LogFatalMissingValue, key)
		return ""
	}
	return value
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("load .env: %v", err)
	}

	cfg := new(Config)

	cfg.HTTP.Host = getEnv(keyHTTPHost)
	cfg.HTTP.Port = getEnv(keyHTTPPort)

	cfg.DB.Host = getEnv(keyDBHost)
	cfg.DB.Port = getEnv(keyDBPort)
	cfg.DB.User = getEnv(keyDBUser)
	cfg.DB.Pass = getEnv(keyDBPass)
	cfg.DB.Name = getEnv(keyDBName)

	return cfg
}
