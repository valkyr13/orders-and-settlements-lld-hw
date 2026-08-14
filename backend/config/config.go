package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	return Config{
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}

	return value
}