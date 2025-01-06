package handlers

import (
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/models"
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/store"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Store  store.Store
	Dbname string
}

func NewHandlers(s store.Store, db string) *Handlers {
	return &Handlers{
		Store:  s,
		Dbname: db,
	}
}

func (h *Handlers) HandlerSet(c *fiber.Ctx) error {
	var req models.JsonReqiest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	h.Store.Set(req.Key, req.Value)
	err := h.Store.Save(h.Dbname)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) HandlerGet(c *fiber.Ctx) error {
	key := c.Query("key")
	if value, ok := h.Store.Get(key); ok {
		return c.JSON(fiber.Map{"value": value})
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "key not found"})
}
