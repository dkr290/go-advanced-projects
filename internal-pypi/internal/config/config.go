package config

import "log"

type AppConfig struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Port     string
	Username string
	Password string
}

func New(username, password, port string) *AppConfig {
	return &AppConfig{
		Username: username,
		Password: password,
		Port:     port,
	}
}
