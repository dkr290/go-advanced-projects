package handlers

// Model handling structures
type PullRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Format string `json:"format"`
}
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Lllam config
type LlamaConfig struct {
	ContextSize int
	GPULayers   int
	NUMA        bool
	Threads     int
	BatchSize   int
	Verbose     bool
}
