package handlers

import (
	"sync"

	"github.com/dkr290/go-advanced-projects/kv-store/pkg/models"
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/store"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Store  store.Store
	Mutext *sync.Mutex
}

func NewHandlers(s store.Store, m *sync.Mutex) *Handlers {
	return &Handlers{
		Store: s,
	}
}

func (h *Handlers) HandlerSet(c *fiber.Ctx) error {
	var req models.JsonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.Store.Set(req.Key, req.Value, req)
	if err != nil {
		return err
	}
	// err = h.Store.Save(req.Database + ".gob")
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).
	// 		JSON(fiber.Map{"error": "Failed to save the data" + err.Error()})
	// }
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) HandlerGet(c *fiber.Ctx) error {
	key := c.Params("key")
	database := c.Params("database")

	if key == "" || database == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Database and key are required"})
	}
	value, exists := h.Store.Get(key, database)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Key not found in database"})
	}

	return c.JSON(fiber.Map{"key": key, "value": value})
}

func (h *Handlers) HandlerGetAllRecords(c *fiber.Ctx) error {
	filename := c.Params("database") + ".jsonl"

	allData, err := h.Store.LoadAll(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load data: " + err.Error(),
		})
	}

	return c.JSON(allData)
}

func (h *Handlers) HandleDelete(c *fiber.Ctx) error {
	key := c.Params("key")
	database := c.Params("database")

	err := h.Store.Delete(key, database)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete the key: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
