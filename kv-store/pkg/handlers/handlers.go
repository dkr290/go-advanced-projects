package handlers

import (
	"sync"

	"github.com/dkr290/go-advanced-projects/kv-store/pkg/models"
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/store"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Store  store.Store
	Mutext sync.Mutex
}

func NewHandlers(s store.Store, db string) *Handlers {
	return &Handlers{
		Store: s,
	}
}

func (h *Handlers) HandlerSet(c *fiber.Ctx) error {
	var req models.JsonReqiest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	h.Mutext.Lock()
	defer h.Mutext.Unlock()

	h.Store.Set(req.Key, req.Value, req)
	err := h.Store.Save(req.Database + ".gob")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Failed to save the data" + err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) HandlerGet(c *fiber.Ctx) error {
	key := c.Query("key")
	database := c.Query("database")
	h.Mutext.Lock()
	defer h.Mutext.Unlock()

	if value, ok := h.Store.Get(key, database); ok {
		return c.JSON(fiber.Map{"value": value})
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "key not found"})
}
