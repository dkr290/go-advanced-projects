package db

import (
	"fmt"
	"os"
	"strconv"
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
	return Config{
		PublicHost: getEnv("HOST", "http://localhost"),
		Port:       getEnv("PORT", "3000"),
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

func getEnv[T int | string](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case int:
		if intVal, err := strconv.Atoi(value); err != nil {
			return any(intVal).(T)
		}
	case string:
		return any(value).(T)
	}
	return defaultValue
}
