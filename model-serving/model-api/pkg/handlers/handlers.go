package handlers

import (
	"fmt"

	"github.com/dkr290/go-advanced-projects/model-api/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	torch "github.com/wangkuiyi/gotorch"
)

type Handlers struct{}

func (h *Handlers) AskHandler(c *fiber.Ctx) error {
	// Parse the request body
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	// Prepare input tensor
	// Tokenize input using Python tokenizer
	tokenizedInput, err := helpers.TokenizeText(req.Prompt, "http://localhost:5001/tokenize")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to prepare input tensor: %v", err),
		})
	}
	// Run inference
	// Convert tokenized input to tensor
	inputTensor := torch.NewTensor(tokenizedInput).To(torch.NewDevice("CPU"))
	model := helpers.LoadModel("model.pt")

	// Return the response
	return c.JSON(Response{
		Answer: answer,
	})
}
