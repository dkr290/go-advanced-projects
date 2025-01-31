package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	Brand     string `json:"brand"`
	Make      string `json:"make"`
	CarModel  string `json:"carmodel"`
	Year      string `json:"year"`
	ImagePath string `json:"imagePath"`
}
