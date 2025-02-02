package handlers

import (
	"github.com/dkr290/go-advanced-projects/cars-htmx/views/home"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) HandleHome(c *fiber.Ctx) error {
	c.Type("html")
	return home.Home().Render(c.Context(), c)
}
