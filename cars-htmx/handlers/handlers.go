package handlers

import (
	"github.com/dkr290/go-advanced-projects/cars-htmx/views/cars"
	"github.com/dkr290/go-advanced-projects/cars-htmx/views/home"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) HandleHome(c *fiber.Ctx) error {
	c.Type("html")
	return home.Home().Render(c.Context(), c)
}

func (h *Handler) HandleListCars(c *fiber.Ctx) error {
	c.Type("html")
	crs, err := h.store.GetAllCars()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retreive cars",
		})
	}
	return cars.CarList(crs).Render(c.Context(), c)
}

func (h *Handler) HandleAddCar(c *fiber.Ctx) error {
	c.Type("html")
	return cars.CarsForm().Render(c.Context(), c)
}
