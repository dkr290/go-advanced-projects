package handlers

// Model handling structures
type PullRequest struct {
	Name   string   `json:"name"`
	URLs   []string `json:"urls"`
	SURL   string   `json:"url"`
	Format string   `json:"format"`
	Meta   struct {
		TotalParts int `json:"total_parts"` // Optional validation
	} `json:"meta"`
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
