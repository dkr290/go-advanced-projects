package sb

import "github.com/nedpals/supabase-go"

func InitDB(sbHost string, sbSecret string) *supabase.Client {
	sbcl := supabase.CreateClient(sbHost, sbSecret)
	return sbcl
}
