package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	Model   ModelConfig
	GPU     GPUConfig
	HF      HuggingFaceConfig
	Python  PythonBackendConfig
	Process ProcessConfig
	Log     LogConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
	Mode string
}

// ModelConfig holds model configuration
type ModelConfig struct {
	Name               string
	Provider           string // huggingface or ollama
	HuggingFaceModelID string
	OllamaModelName    string
	CacheDir           string
	MaxFrames          int
	DefaultFPS         int
	DefaultWidth       int
	DefaultHeight      int
	MaxVideoDuration   int
}

// GPUConfig holds GPU configuration
type GPUConfig struct {
	Enabled        bool
	DeviceID       int
	MemoryFraction float64
}

// HuggingFaceConfig holds Hugging Face configuration
type HuggingFaceConfig struct {
	Token  string
	APIURL string
}

// PythonBackendConfig holds Python backend configuration
type PythonBackendConfig struct {
	URL     string
	Enabled bool
	Timeout time.Duration
}

// ProcessConfig holds processing configuration
type ProcessConfig struct {
	MaxConcurrentRequests int
	RequestTimeout        time.Duration
	UploadMaxSize         int64
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// LoadConfig loads configuration from environment and config files
func LoadConfig() (*Config, error) {
	// Set defaults
	setDefaults()

	config := &Config{
		Server: ServerConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetInt("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Model: ModelConfig{
			Name:               viper.GetString("MODEL_NAME"),
			Provider:           getProvider(),
			HuggingFaceModelID: viper.GetString("HUGGINGFACE_MODEL_ID"),
			OllamaModelName:    viper.GetString("OLLAMA_MODEL_NAME"),
			CacheDir:           viper.GetString("MODEL_CACHE_DIR"),
			MaxFrames:          viper.GetInt("MAX_FRAMES"),
			DefaultFPS:         viper.GetInt("DEFAULT_FPS"),
			DefaultWidth:       viper.GetInt("DEFAULT_WIDTH"),
			DefaultHeight:      viper.GetInt("DEFAULT_HEIGHT"),
			MaxVideoDuration:   viper.GetInt("MAX_VIDEO_DURATION"),
		},
		GPU: GPUConfig{
			Enabled:        viper.GetBool("ENABLE_GPU"),
			DeviceID:       viper.GetInt("GPU_DEVICE_ID"),
			MemoryFraction: viper.GetFloat64("GPU_MEMORY_FRACTION"),
		},
		HF: HuggingFaceConfig{
			Token:  viper.GetString("HUGGINGFACE_TOKEN"),
			APIURL: viper.GetString("HUGGINGFACE_API_URL"),
		},
		Python: PythonBackendConfig{
			URL:     viper.GetString("PYTHON_BACKEND_URL"),
			Enabled: viper.GetBool("PYTHON_BACKEND_ENABLED"),
			Timeout: time.Duration(viper.GetInt("REQUEST_TIMEOUT")) * time.Second,
		},
		Process: ProcessConfig{
			MaxConcurrentRequests: viper.GetInt("MAX_CONCURRENT_REQUESTS"),
			RequestTimeout:        time.Duration(viper.GetInt("REQUEST_TIMEOUT")) * time.Second,
			UploadMaxSize:         parseSize(viper.GetString("UPLOAD_MAX_SIZE")),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_MODE", "release")

	viper.SetDefault("MODEL_NAME", "Wan2.1")
	viper.SetDefault("HUGGINGFACE_MODEL_ID", "Lightricks/LTX-Video")
	viper.SetDefault("MODEL_CACHE_DIR", "./models")
	viper.SetDefault("USE_HUGGINGFACE", true)
	viper.SetDefault("USE_OLLAMA", false)

	viper.SetDefault("ENABLE_GPU", true)
	viper.SetDefault("GPU_DEVICE_ID", 0)
	viper.SetDefault("GPU_MEMORY_FRACTION", 0.9)

	viper.SetDefault("HUGGINGFACE_API_URL", "https://huggingface.co")

	viper.SetDefault("PYTHON_BACKEND_URL", "http://localhost:5000")
	viper.SetDefault("PYTHON_BACKEND_ENABLED", true)

	viper.SetDefault("MAX_FRAMES", 128)
	viper.SetDefault("DEFAULT_FPS", 24)
	viper.SetDefault("DEFAULT_WIDTH", 512)
	viper.SetDefault("DEFAULT_HEIGHT", 512)
	viper.SetDefault("MAX_VIDEO_DURATION", 10)

	viper.SetDefault("MAX_CONCURRENT_REQUESTS", 2)
	viper.SetDefault("REQUEST_TIMEOUT", 300)
	viper.SetDefault("UPLOAD_MAX_SIZE", "100MB")

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
}

func getProvider() string {
	if viper.GetBool("USE_HUGGINGFACE") {
		return "huggingface"
	}
	if viper.GetBool("USE_OLLAMA") {
		return "ollama"
	}
	return "huggingface"
}

func parseSize(sizeStr string) int64 {
	// Simple parser for sizes like "100MB"
	var size int64
	var unit string
	fmt.Sscanf(sizeStr, "%d%s", &size, &unit)

	switch unit {
	case "KB", "kb":
		return size * 1024
	case "MB", "mb":
		return size * 1024 * 1024
	case "GB", "gb":
		return size * 1024 * 1024 * 1024
	default:
		return size
	}
}
