package main

import "log"

type config struct {
	addr      string
	staticDir string
	dsn       string
}
type appconfig struct {
	errotLog *log.Logger
	infoLog  *log.Logger
}
