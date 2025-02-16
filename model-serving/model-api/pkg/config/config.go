package config

// Lllam config
type LlamaConfig struct {
	ContextSize int
	GPULayers   int
	NUMA        bool
	Threads     int
	BatchSize   int
	Verbose     bool
}
