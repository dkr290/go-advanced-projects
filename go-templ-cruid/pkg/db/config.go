package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the env file")
	}
	return Config{
		PublicHost: getEnv("HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "Password123"),
		DBAddress: fmt.Sprintf(
			"%s:%s",
			getEnv("DB_HOST", "mysql"),
			getEnv("DB_PORT", "3306"),
		),
		DBName: getEnv("DB_NAME", "todo"),
	}
}

func getEnv[T int | string](key string, fallback T) T {
	if value, ok := os.LookupEnv(key); ok {
		var result T
		switch any(fallback).(type) {
		case int:
			if i, err := strconv.Atoi(value); err == nil {
				result = any(i).(T)
			} else {
				result = fallback
			}
		case string:
			result = any(value).(T)
		default:
			result = fallback
		}
		return result
	}

	return fallback
}
