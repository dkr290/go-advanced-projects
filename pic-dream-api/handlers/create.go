package handlers

import (
	"github.com/nedpals/supabase-go"
	"github.com/uptrace/bun"
)

type Handlers struct {
	sb                  supabase.Client
	github_redirect_url string
	Bun                 *bun.DB
}

func NewHandlers(sb supabase.Client, githubhRdr string, bundb bun.DB) *Handlers {
	return &Handlers{
		sb:                  sb,
		github_redirect_url: githubhRdr,
		Bun:                 &bundb,
	}
}
