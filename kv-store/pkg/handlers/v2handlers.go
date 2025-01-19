package handlers

import (
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/models"
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/store"
	"github.com/gofiber/fiber/v2"
)

type V2Handlers struct {
	V2Store store.V2Store
}

func NewV2Handlers(s store.V2Store) *V2Handlers {
	return &V2Handlers{
		V2Store: s,
	}
}

func (h *V2Handlers) V2HandlerSet(c *fiber.Ctx) error {
	var req models.V2JsonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.V2Store.Set(req.Key, req.Value, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *V2Handlers) V2HandlerGet(c *fiber.Ctx) error {
	key := c.Params("key")
	database := c.Params("database")

	if key == "" || database == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Database and key are required"})
	}
	value, exists := h.V2Store.Get(key, database)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Key not found in database"})
	}

	return c.JSON(fiber.Map{"key": key, "value": value})
}

func (h *V2Handlers) V2HandlerGetAllRecords(c *fiber.Ctx) error {
	filename := c.Params("database") + ".jsonl"

	allData, err := h.V2Store.LoadAll(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load data: " + err.Error(),
		})
	}

	return c.JSON(allData)
}

func (h *V2Handlers) V2HandleDelete(c *fiber.Ctx) error {
	key := c.Params("key")
	database := c.Params("database")

	err := h.V2Store.Delete(key, database)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete the key: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
