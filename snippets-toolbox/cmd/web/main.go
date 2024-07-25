package main

import (
	"log"
	"net/http"
)

const PORT = ":4000"

func main() {

	mux := http.NewServeMux()
	//create a file server which serves files out of "./ui/static direct all "
	//path is relative to the project directory root
	fs := http.FileServer(http.Dir("./ui/static/"))

	//use the handler function to register the fileserver as handler
	mux.Handle("/static", http.StripPrefix("/static", fs))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {

		log.Fatal(err)
	}
}
