package main

import "os"

const (
	packageDir = "./packages"
)

func getCredentials() (username string, password string) {
	username = os.Getenv("USERNAME")
	if len(username) == 0 {
		username = "admin"
	}
	password = os.Getenv("PASSWORD")
	if len(password) == 0 {
		password = "password"
	}
	return
}
