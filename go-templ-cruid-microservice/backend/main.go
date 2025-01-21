package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/handlers"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/pkg/db"
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

	// get all tasks and return json for to be consumed by frontend
	r.Get("/tasks", h.HandleGetAllTasks)
	// add task post handler for adding task and return sucess json response to be consumed by frontend
	r.Post("/addTask", h.HandleAddTask)

	// update the task
	r.Put("/task/{id}", h.HandleUpdateTask)

	// delete the task
	r.Delete("/task/{id}", h.HandleDeleteTask)

	r.Get("/task/{id}", h.HandleGetTaskById)

	port := os.Getenv("HTTP_LISTEN_ADDR")
	if port == "" {
		port = "localhost:3000"
	}
	slog.Info("application is running", "port", port)
	return http.ListenAndServe(port, r)
}
