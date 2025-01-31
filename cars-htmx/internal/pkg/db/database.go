package db

import (
	"fmt"
	"time"

	"github.com/dkr290/go-advanced-projects/cars-htmx/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqlLiteDb(config DbConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(config.DBName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Car{})
	if err != nil {
		return nil, err
	}
	return
}

// in case of Postgresql
func InitPostgresDb(config DbConfig, numRetries int) (db *gorm.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.DBUser, config.DBPassword, config.DBName)

	for i := 0; i <= numRetries; i++ {
		db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

		if i > numRetries {
			return nil, fmt.Errorf("Failed to connect to the database %v", err)
		}
		if err != nil {
			fmt.Printf("Trying to connect to the database %d time \n", i)
			fmt.Println(err)
			time.Sleep(2 * time.Second)
		} else {
			fmt.Println("Connected to the database")
			break
		}
	}
	err = db.AutoMigrate(&models.Car{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
