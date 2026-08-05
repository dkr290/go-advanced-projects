// Package utils contains reusable domain logic shared across handlers.
package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SearchResult represents a single DuckDuckGo search result.
type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// DuckDuckGoClient performs web searches against DuckDuckGo's HTML endpoint.
// The HTML endpoint (html.duckduckgo.com/html/) is used because it does not
// require an API key, unlike the official Instant Answer API which only
// returns limited structured data.
type DuckDuckGoClient struct {
	httpClient *http.Client
	userAgent  string
	baseURL    string
}

// NewDuckDuckGoClient builds a client using the provided http.Client and
// User-Agent header value.
func NewDuckDuckGoClient(httpClient *http.Client, userAgent string) *DuckDuckGoClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &DuckDuckGoClient{
		httpClient: httpClient,
		userAgent:  userAgent,
		baseURL:    "https://html.duckduckgo.com/html/",
	}
}

// Search queries DuckDuckGo for the given topic and returns up to maxResults
// results. The context controls cancellation/timeout of the outbound HTTP
// request.
func (c *DuckDuckGoClient) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query must not be empty")
	}
	if maxResults <= 0 {
		maxResults = 10
	}

	form := url.Values{}
	form.Set("q", query)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse response body: %w", err)
	}

	return parseResults(doc, maxResults), nil
}

func parseResults(doc *goquery.Document, maxResults int) []SearchResult {
	results := make([]SearchResult, 0, maxResults)

	doc.Find(".result").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if len(results) >= maxResults {
			return false
		}

		linkSel := s.Find("a.result__a").First()
		title := strings.TrimSpace(linkSel.Text())
		href, exists := linkSel.Attr("href")
		if title == "" || !exists {
			return true
		}

		snippet := strings.TrimSpace(s.Find(".result__snippet").First().Text())

		results = append(results, SearchResult{
			Title:   title,
			URL:     cleanDuckDuckGoRedirect(href),
			Snippet: snippet,
		})
		return true
	})

	return results
}

// cleanDuckDuckGoRedirect unwraps DuckDuckGo's redirect links
// (//duckduckgo.com/l/?uddg=<encoded-url>&...) into the plain target URL.
func cleanDuckDuckGoRedirect(href string) string {
	if strings.HasPrefix(href, "//") {
		href = "https:" + href
	}

	u, err := url.Parse(href)
	if err != nil {
		return href
	}
	if strings.Contains(u.Host, "duckduckgo.com") && u.Path == "/l/" {
		if target := u.Query().Get("uddg"); target != "" {
			if decoded, err := url.QueryUnescape(target); err == nil {
				return decoded
			}
		}
	}
	return href
}
