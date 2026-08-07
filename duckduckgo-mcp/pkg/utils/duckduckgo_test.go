package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)


const sampleHTML = `
<html><body>
<div class="result">
	<a class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fgolang.org%2F&amp;rut=abc">Go Programming Language</a>
	<a class="result__snippet">The Go programming language.</a>
</div>
<div class="result">
	<a class="result__a" href="https://example.com/page">Example Page</a>
	<a class="result__snippet">An example page.</a>
</div>
</body></html>`

func TestDuckDuckGoClient_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if q := r.FormValue("q"); q != "golang" {
			t.Errorf("expected q=golang, got %q", q)
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(sampleHTML))
	}))
	defer server.Close()

	client := NewDuckDuckGoClient(server.Client(), "test-agent")
	client.baseURL = server.URL

	results, err := client.Search(context.Background(), "golang", 10)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	want := []SearchResult{
		{Title: "Go Programming Language", URL: "https://golang.org/", Snippet: "The Go programming language."},
		{Title: "Example Page", URL: "https://example.com/page", Snippet: "An example page."},
	}

	for i, r := range results {
		if r != want[i] {
			t.Errorf("result[%d] = %+v, want %+v", i, r, want[i])
		}
	}
}

func TestDuckDuckGoClient_Search_EmptyQuery(t *testing.T) {
	client := NewDuckDuckGoClient(http.DefaultClient, "test-agent")
	if _, err := client.Search(context.Background(), "   ", 10); err == nil {
		t.Fatal("expected error for empty query, got nil")
	}
}

func TestDuckDuckGoClient_Search_MaxResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(sampleHTML))
	}))
	defer server.Close()

	client := NewDuckDuckGoClient(server.Client(), "test-agent")
	client.baseURL = server.URL

	results, err := client.Search(context.Background(), "golang", 1)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestDuckDuckGoClient_Search_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewDuckDuckGoClient(server.Client(), "test-agent")
	client.baseURL = server.URL

	if _, err := client.Search(context.Background(), "golang", 10); err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestCleanDuckDuckGoRedirect(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "protocol relative redirect with uddg",
			in:   "//duckduckgo.com/l/?uddg=https%3A%2F%2Fgolang.org%2F&rut=abc",
			want: "https://golang.org/",
		},
		{
			name: "plain url passthrough",
			in:   "https://example.com/page",
			want: "https://example.com/page",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanDuckDuckGoRedirect(tc.in)
			if got != tc.want {
				t.Fatalf("cleanDuckDuckGoRedirect(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
