package handlers

import "github.com/nedpals/supabase-go"

type Handlers struct {
	sb supabase.Client
}

func NewHandlers(sb supabase.Client) *Handlers {
	return &Handlers{

		sb: sb,
	}
}
