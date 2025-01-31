package main

import (
	"log"

	"github.com/dkr290/go-advanced-projects/cars-htmx/internal/pkg/db"
)

func main() {
	conf := db.InitConfig()
	d, err := db.InitSqlLiteDb(conf)
	if err != nil {
		log.Fatal(err)
	}
	database := db.Storage{
		Db: d,
	}
	//  just to use database

	log.Println("Database name:", database.Db.Name())
}
