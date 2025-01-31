package db

import (
	"os"
	"strconv"
)

type DbConfig struct {
	Host       string
	Port       int
	DBUser     string
	DBPassword string
	DBName     string
}

func InitConfig() DbConfig {
	return DbConfig{
		Host:       getEnv("HOST", "localhost"),
		Port:       getEnv("PORT", 5432),
		DBUser:     getEnv("DBUSER", "postgres"),
		DBPassword: getEnv("DBPASS", "password"),
		DBName:     getEnv("DB", "postgres"),
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
