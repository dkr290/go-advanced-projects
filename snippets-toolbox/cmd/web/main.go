package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	var cfg config

	//define a new command-line flag with the name addr and default value
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any
	// errors are
	// encountered during parsing the application will be terminated.

	flag.Parse()

	mux := http.NewServeMux()
	//create a file server which serves files out of "./ui/static direct all "
	//path is relative to the project directory root
	fs := http.FileServer(http.Dir("./ui/static/"))

	//use the handler function to register the fileserver as handler
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	// The value returned from the flag.String() function is a pointer to theflag
	// value, not the value itself. So we need to dereference the pointer
	// prefix it with the * symbol) before using it. Note that it is using
	// log.Printf() function to interpolate the address with the log message.

	log.Printf("Starting server on %s", cfg.addr)
	if err := http.ListenAndServe(cfg.addr, mux); err != nil {

		log.Fatal(err)
	}
}
