package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const Port = ":4000"

func home(w http.ResponseWriter, r *http.Request) {
	//check if the current path / exactly match and if not use http.NotFound()

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, _ = w.Write([]byte("Hello from snippets-toolbox"))

}

// snippetView() handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	//extract the value of the id parameter from the query string
	// convert it to integer and compare if it is less then 1 return 404 page
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

//snippetCreate() handler to create snippets

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// use r.Method to check if there is post request or not and if not status 405 to be displayed
	if r.Method != http.MethodPost { // use constant
		w.Header().Set("Allow", "POST") // it can be the string also
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte("Create new snippet..."))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")

	if err := http.ListenAndServe(Port, mux); err != nil {
		log.Fatal(err)
	}
}
