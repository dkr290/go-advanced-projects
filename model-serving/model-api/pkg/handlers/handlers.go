package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	ModelsDir string
}

func NewHandlers(modelsDir string) *Handlers {
	return &Handlers{
		ModelsDir: modelsDir,
	}
}

// Handlers
func (h *Handlers) PullModel(c *fiber.Ctx) error {
	var req PullRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	modelPath := filepath.Join(h.ModelsDir, req.Name)
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
