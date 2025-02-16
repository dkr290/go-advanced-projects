package handlers

// Model handling structures
type PullRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Format string `json:"format"`
}
type GenerateRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"tetmperature"`
	TopP        float32 `json:"topp"`
	TopK        int     `json:"topk"`
	MaxTokens   int     `json:"maxtoklens"`
	Seed        int     `json:"seed"`
}
