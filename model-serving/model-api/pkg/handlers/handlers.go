package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dkr290/go-advanced-projects/model-serving/model-api/pkg/config"
	llama "github.com/go-skynet/go-llama.cpp"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	ModelsDir   string
	Sem         chan struct{}
	LlamaConfig *config.LlamaConfig
}

func NewHandlers(modelsDir string, sem chan struct{}, llamaConfig *config.LlamaConfig) *Handlers {
	return &Handlers{
		ModelsDir:   modelsDir,
		Sem:         sem,
		LlamaConfig: llamaConfig,
	}
}

// Handlers
func (h *Handlers) PullModelgguf(c *fiber.Ctx) error {
	var req PullRequest // The pull request for download model file in GGUF

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	modelPath := filepath.Join(h.ModelsDir, req.Name)

	// check if the models is already downloaded
	if _, err := os.Stat(modelPath); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "model already exists"})
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to download model",
		})
	}

	outFile, err := os.Create(modelPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handlers) GenerateRequest(c *fiber.Ctx) error {
	select {
	case h.Sem <- struct{}{}:
		defer func() { <-h.Sem }()
	default:
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "server busy, try again later",
		})
	}

	var req GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	modelPath := filepath.Join(h.ModelsDir, req.Model)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "model not found"})
	}

	l, err := llama.New(
		modelPath,
		llama.SetContext(h.LlamaConfig.ContextSize),
		llama.SetGPULayers(h.LlamaConfig.GPULayers),
		llama.EnableF16Memory,
		llama.EnableEmbeddings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to initialize model: %s", err.Error()),
		})
	}
	defer l.Free()

	// Set inference parameters from request
	opts := []llama.PredictOption{
		llama.SetTemperature(req.Temperature),
		llama.SetTopP(req.TopP),
		llama.SetTopK(req.TopK),
		llama.SetTokens(req.MaxTokens),
		llama.SetThreads(h.LlamaConfig.Threads),
		llama.SetSeed(req.Seed),
	}
	// Run prediction with context
	result, err := l.Predict(req.Prompt, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("prediction failed: %s", err.Error()),
		})
	}

	return c.JSON(fiber.Map{
		"response": result,
	})
}

func (h *Handlers) ListModels(c *fiber.Ctx) error {
	files, err := os.ReadDir(h.ModelsDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	models := make([]fiber.Map, 0)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		models = append(models, fiber.Map{
			"name": file.Name(),
			"size": info.Size(),
		})
	}

	return c.JSON(models)
}

func (h *Handlers) DeleteModel(c *fiber.Ctx) error {
	modelName := c.Params("name")
	modelPath := filepath.Join(h.ModelsDir, modelName)

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "model not found"})
	}

	if err := os.Remove(modelPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
