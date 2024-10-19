package handlers

import "github.com/nedpals/supabase-go"

type Handlers struct {
	sb                  supabase.Client
	github_redirect_url string
}

func NewHandlers(sb supabase.Client, githubhRdr string) *Handlers {
	return &Handlers{
		sb:                  sb,
		github_redirect_url: githubhRdr,
	}
}
