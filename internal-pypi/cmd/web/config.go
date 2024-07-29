package main

import (
	"os"
	"sort"
	"strings"
)

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
func sortPackages(packages []string) []string {
	sort.Slice(packages, func(i, j int) bool {
		// Split the package paths into components
		partsI := strings.Split(packages[i], "/")
		partsJ := strings.Split(packages[j], "/")

		// Compare the package names (last part of the path)
		nameI := partsI[len(partsI)-1]
		nameJ := partsJ[len(partsJ)-1]

		return nameI < nameJ
	})

	return packages
}
