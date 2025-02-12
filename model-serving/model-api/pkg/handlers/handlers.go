package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/dkr290/go-advanced-projects/model-serving/model-api/pkg/helpers"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	ModelsDir    string
	Sem          chan struct{}
	LlamaCppPath string
}

func NewHandlers(modelsDir string, sem chan struct{}, llamacpppath string) *Handlers {
	return &Handlers{
		ModelsDir:    modelsDir,
		Sem:          sem,
		LlamaCppPath: llamacpppath,
	}
}

// Handlers
func (h *Handlers) PullModelgguf(c *fiber.Ctx) error {
	var req PullRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	modelPath := filepath.Join(h.ModelsDir, req.Name)
	if _, err := os.Stat(modelPath); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "model already exists"})
	}

	resp, err := http.Get(req.SURL)
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

func (h *Handlers) PullSafeTensors(c *fiber.Ctx) error {
	var req PullRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	modelDir := filepath.Join(h.ModelsDir, req.Name)
	tempDir := filepath.Join(h.ModelsDir, "temp_"+req.Name)

	// Create directories
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Download all parts concurrently
	errChan := make(chan error, len(req.URLs))
	var wg sync.WaitGroup

	for i, url := range req.URLs {
		wg.Add(1)
		go func(idx int, downloadURL string) {
			defer wg.Done()

			resp, err := http.Get(downloadURL)
			if err != nil {
				errChan <- fmt.Errorf("part %d download failed: %w", idx, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("part %d bad status: %s", idx, resp.Status)
				return
			}
			partName := fmt.Sprintf("part-%04d.safetensors", idx)
			outPath := filepath.Join(tempDir, partName)

			outFile, err := os.Create(outPath)
			if err != nil {
				errChan <- fmt.Errorf("part %d create failed: %w", idx, err)
				return
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, resp.Body); err != nil {
				errChan <- fmt.Errorf("part %d write failed: %w", idx, err)
				return
			}
		}(i, url)
	}

	wg.Wait()
	close(errChan)
	// Check for errors
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "partial download",
			"details": errors,
		})
	}

	// Conversion logic for multi-file safetensors
	if req.Format == "safetensors-multi" {
		if err := helpers.ConvertMultiSafetensors(tempDir, modelDir); err != nil {
			os.RemoveAll(tempDir)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("conversion failed: %v", err),
			})
		}
	}

	// Cleanup
	os.RemoveAll(tempDir)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, h.LlamaCppPath, "-m", modelPath, "-p", req.Prompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  err.Error(),
			"stderr": stderr.String(),
		})
	}

	return c.JSON(fiber.Map{
		"response": stdout.String(),
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
