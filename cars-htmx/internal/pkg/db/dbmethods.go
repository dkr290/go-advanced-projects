package db

import (
	"github.com/dkr290/go-advanced-projects/cars-htmx/internal/models"
	"gorm.io/gorm"
)

type Database interface {
	InsertCar(c *models.CarPostRequest) error
	GetAllCars() ([]models.Car, error)
	DeleteCar(id string) error
	FindCarsByNameMakeOrBrand(search string) ([]models.Car, error)
}

type Storage struct {
	Db *gorm.DB
}

// dbmethods
func (s *Storage) InsertCar(c *models.CarPostRequest) error {
	tx := s.Db.Create(&models.Car{
		Brand:     c.Brand,
		Make:      c.Make,
		CarModel:  c.Model,
		Year:      c.Year,
		ImagePath: c.ImagePath,
	})

	return tx.Error
}

func (p *Storage) GetAllCars() ([]models.Car, error) {
	var c []models.Car
	tx := p.Db.Find(&c)
	return c, tx.Error
}

func (p *Storage) DeleteCar(id string) error {
	tx := p.Db.Delete(&id)
	return tx.Error
}

func (p *Storage) FindCarsByNameMakeOrBrand(search string) ([]models.Car, error) {
	var cars []models.Car
	query := "%" + search + "%" // Wildcard for like search
	err := p.Db.Where("brand LIKE ? OR carmodel LIKE ? OR make LIKE ?", query, query, query).
		Find(&cars).
		Error
	if err != nil {
		return nil, err
	}

	return cars, nil
}
