package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dkr290/go-advanced-projects/ecom/config"
	"github.com/dkr290/go-advanced-projects/ecom/db"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		upFlag   bool
		downFlag bool
	)

	flag.BoolVar(&upFlag, "up", false, "Make migrations")
	flag.BoolVar(&downFlag, "down", false, "Make migrations")
	flag.Parse()

	mdb := db.MysqlDB{}

	db, err := mdb.InitDB(mysqlCfg.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	if upFlag {

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("make migrations up")
	} else if downFlag {

		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("make migrations down")
	} else {
		fmt.Println("Please specofy either -up or -down")
	}

}
