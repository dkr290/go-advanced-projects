package main

import (
	"log"

	"github.com/dkr290/go-advanced-projects/ecom/cmd/api"
	"github.com/dkr290/go-advanced-projects/ecom/config"
	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/go-sql-driver/mysql"
)

func main() {

	mdb := db.MysqlDB{}

	err := mdb.InitDB(mysql.Config{
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

	server := api.New(":8080", &mdb)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
