package handlers

import (
	"github.com/dkr290/go-advanced-projects/pic-dream-api/pkg/db"
	"github.com/nedpals/supabase-go"
)

type Handlers struct {
	sb                  supabase.Client
	github_redirect_url string
	Bun                 db.PictureDatabase
}

func NewHandlers(sb supabase.Client, githubhRdr string, bundb db.PictureDatabase) *Handlers {
	return &Handlers{
		sb:                  sb,
		github_redirect_url: githubhRdr,
		Bun:                 bundb,
	}
}
