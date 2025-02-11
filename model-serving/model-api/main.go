package main

const (
	modelsDir      = "models"
	llamaCPPPath   = "llama.cpp/build/bin/llama-cli" // Update this path
	maxConcurrency = 4
)

var sem = make(chan struct{}, maxConcurrency)

func main() {
}
