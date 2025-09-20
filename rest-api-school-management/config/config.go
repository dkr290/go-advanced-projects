package config

import (
	"flag"
	"os"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBDatabase string
	DBPort     string
	DBHost     string
}

func LoadConfig() *Config {
	c := &Config{}
	c.GetFlags()
	return c
}

func (c *Config) GetFlags() {
	flag.StringVar(
		&c.Port,
		"app-port",
		"8082",
		"Application port",
	)
	flag.StringVar(&c.DBUser, "database-user", "root", "Database user")
	flag.StringVar(&c.DBPassword, "database-password", "password", "Database password")
	flag.StringVar(&c.DBDatabase, "database", "school_test", "Database to use")
	flag.StringVar(&c.DBHost, "database-host", "mariadb", "The database host to use or ip")
	flag.StringVar(
		&c.DBPort,
		"database-port",
		"3306",
		"Database port for Mariadb database connection",
	)
	flag.Parse()

	if dbhost := getEnv("DATABASE_HOST"); dbhost != "" {
		c.DBHost = dbhost
	}

	if port := getEnv("PORT"); port != "" {
		c.Port = port
	}
	if dbuser := getEnv("DB_USER"); dbuser != "" {
		c.DBUser = dbuser
	}

	if dbpass := getEnv("DB_PASSWORD"); dbpass != "" {
		c.DBPassword = dbpass
	}
	if database := getEnv("DATABASE"); database != "" {
		c.DBDatabase = database
	}
	if dbport := getEnv("DATABASE_PORT"); dbport != "" {
		c.DBPort = dbport
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
