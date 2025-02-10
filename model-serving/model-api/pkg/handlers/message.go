package handlers

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Answer string `json:"Answer"`
}
