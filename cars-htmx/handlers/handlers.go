package handlers

import (
	"log"

	"github.com/dkr290/go-advanced-projects/cars-htmx/internal/models"
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
	log.Println("hit add car")

	params := models.CarPostRequest{
		Model:     c.FormValue("model"),
		Brand:     c.FormValue("brand"),
		Make:      c.FormValue("make"),
		Year:      c.FormValue("Year"),
		ImagePath: c.FormValue("imagepath"),
	}
	err := h.store.InsertCar(&params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error insert car",
		})
	}
	return c.Redirect("/cars")
}

func (h *Handler) ShowCarForm(c *fiber.Ctx) error {
	return cars.CarsForm().Render(c.Context(), c)
}
