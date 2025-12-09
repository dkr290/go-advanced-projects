package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Movie struct {
	Identifier  string `json:"identifier"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Response struct {
	Response struct {
		Docs []Movie `json:"docs"`
	} `json:"response"`
}

func searchMovies(keyword string, rows int) error {
	basURL := "https://archive.org/advancedsearch.php"
	params := url.Values{}
	params.Set("q", fmt.Sprintf("mediaType:movies AND (%s)", keyword))
	params.Add("fl[]", "identifier")
	params.Add("fl[]", "title")
	params.Add("fl[]", "description")
	params.Set("rows", fmt.Sprintf("%d", rows))
	params.Set("output", "json")

	resp, err := http.Get(basURL + "?" + params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	for _, movie := range result.Response.Docs {
		fmt.Printf(
			"Title: %s\nIdentifier: %s\n Description %.200s...\n\n",
			movie.Title,
			movie.Identifier,
			movie.Description,
		)
	}
	return nil
}

func main() {
	k := flag.String("keyword", "", "keyword")
	flag.Parse()
	if *k == "" {
		fmt.Println("keyword is empty ")
		os.Exit(1)
	}
	if err := searchMovies(*k, 5); err != nil {
		fmt.Println("error:", err)
	}
}
