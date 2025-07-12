package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App        AppConfig
	HTTPServer HTTPConfig
	Database   DatabaseConfig
	Logging    LoggerConfig
	Services   Services
	JWTSecret  string `mapstructure:"JWT_SECRET"`
}

type AppConfig struct {
	Version string
}

type HTTPConfig struct {
	Port              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type LoggerConfig struct {
	Enable     bool
	LogsDir    string
	Level      string
	Format     string
	SavingDays int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type Services struct {
	MobileApp Service
}

type Service struct {
	Host string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %s", err.Error())
	}

	return &Config{
		App: AppConfig{
			Version: getEnv("VERSION", "1.0.0"),
		},
		HTTPServer: HTTPConfig{
			Port:              getEnv("SERVER_PORT", "6004"),
			ReadTimeout:       time.Second * 10,
			ReadHeaderTimeout: time.Second * 20,
			WriteTimeout:      time.Second * 20,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "192.168.29.138"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "db_name"),
		},
		Logging: LoggerConfig{
			Enable:     getEnvAsBool("LOGGER_ENABLE", true),
			LogsDir:    getEnv("LOGGER_LOGS_DIR", "./logs"),
			Level:      getEnv("LOGGER_LOG_LEVEL", "DEBUG"),
			SavingDays: getEnvAsInt("LOGGER_SAVING_DAYS", 5),
		},
		Services: Services{
			MobileApp: Service{
				Host: getEnv("API_URL", "http://localhost:8080"),
			},
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	val, _ := strconv.ParseBool(value)
	return val
}
