package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid/handlers"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/helpers"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/pkg/db"
	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
)

func main() {
	database, err := db.InitDB(mysql.Config{
		User:                 db.Envs.DBUser,
		Passwd:               db.Envs.DBPassword,
		Addr:                 db.Envs.DBAddress,
		DBName:               db.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal("Fatal error loading environment for mysql connection", err)
		return
	}

	mdb := db.MysqlDatabase{
		DB: database,
	}

	log.Fatal(Run(mdb))
}

func Run(mdb db.MysqlDatabase) error {
	h := handlers.NewHandlers(&mdb)

	r := chi.NewRouter()
	r.Get("/", helpers.MakeHandler(h.HandleHome))
	// get all tasks
	r.Get("/tasks", helpers.MakeHandler(h.HandleFetchTasks))
	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application is running", "port", port)
	return http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), r)
}
