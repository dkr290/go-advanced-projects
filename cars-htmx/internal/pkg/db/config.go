package db

import (
	"flag"
	"fmt"
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
	var dbName string
	flag.StringVar(&dbName, "db", "sqllite", "Database name postgres or sqllite")
	// Parse flags
	flag.Parse()
	// Check if flag was explicitly set
	if dbName == "" {
		fmt.Fprintln(os.Stderr, "\nError: -db flag is required")
		flag.Usage() // Show help message
		os.Exit(1)
	}

	// Validate allowed values
	if dbName != "postgres" && dbName != "sqllite" {
		fmt.Fprintf(os.Stderr, "Error: Invalid DB type '%s' (allowed: postgres, sqllite)\n", dbName)
		os.Exit(1)
	}

	return DbConfig{
		Host:       getEnv("HOST", "localhost"),
		Port:       getEnv("PORT", 5432),
		DBUser:     getEnv("DBUSER", "postgres"),
		DBPassword: getEnv("DBPASS", "password"),
		DBName:     dbName,
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
