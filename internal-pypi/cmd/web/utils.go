package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"sort"
	"strings"
)

const (
	packageDir = "./packages"
)

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
func basicAuth(next http.HandlerFunc, username string, password string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
func (app *Config) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *Config) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Config) clientLog(s string, arg ...any) {
	app.InfoLog.Printf(s, arg...)
}
func (app *Config) errorLog(s string, arg ...any) {
	app.ErrorLog.Printf(s, arg...)
}
