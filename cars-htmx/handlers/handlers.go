package handlers

import (
	"strconv"

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

	params := models.CarPostRequest{
		Model:     c.FormValue("model"),
		Brand:     c.FormValue("brand"),
		Make:      c.FormValue("make"),
		Year:      c.FormValue("year"),
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

func (h *Handler) HandleDeleteCar(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error invalid id",
		})
	}
	err = h.store.DeleteCar(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Error item not found",
		})
	}
	// Fetch the updated list of cars
	newCars, err := h.store.GetAllCars()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cars",
		})
	}
	return cars.CarListWithToast(newCars, "Car deleted sucessfully").Render(c.Context(), c)
}
