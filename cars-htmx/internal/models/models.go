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

type CarPostRequest struct {
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Make      string `json:"make"`
	Year      string `json:"year"`
	ImagePath string `json:"imagepath"`
}
