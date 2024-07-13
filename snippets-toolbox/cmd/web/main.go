package main

import (
	"log"
	"net/http"
)

const PORT = ":4000"

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {

		log.Fatal(err)
	}
}
