package main

import (
	"log"
	"net/http"
	"os"
)

const PORT = ":4000"

func main() {

	mux := http.NewServeMux()

	if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
		log.Fatalf("could not create upload directory: %v", err)
	}
	//create a file server which serves files out of "./ui/static direct all "
	//path is relative to the project directory root
	fs := http.FileServer(http.Dir("./static/"))

	//use the handler function to register the fileserver as handler
	mux.Handle("/static", http.StripPrefix("/static", fs))
	mux.HandleFunc("/", basicAuth(indexHandler))
	mux.HandleFunc("/simple/", basicAuth(simpleHandler))
	mux.HandleFunc("/packages/", basicAuth(packageHandler))
	mux.HandleFunc("/upload", basicAuth(uploadHandler))
	mux.HandleFunc("/about", basicAuth(aboutHandler))
	mux.HandleFunc("/contact", basicAuth(contactHandler))

	log.Print("Starting server on", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {

		log.Fatal(err)
	}
}
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	username, password := getCredentials()

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
