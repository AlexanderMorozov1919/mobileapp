package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading .env: %v", err)
	}
	return nil
}

func LoadDBConfig() (DBConfig, error) {
	required := []string{"DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, key := range required {
		if os.Getenv(key) == "" {
			return DBConfig{}, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "med_data"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}, nil
}

func InitDB() (*gorm.DB, error) {
	if err := LoadEnv(); err != nil {
		return nil, fmt.Errorf("error loading environment: %w", err)
	}

	cfg, err := LoadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading DB config: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Successfully connected to database")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
