package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/handlers"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/pkg/db"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/helpers"
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

	// fetch update form for update by ID
	r.Get("/gettaskupdateform/{id}", helpers.MakeHandler(h.HandleGetTaskUpdateForm))

	// update the task
	r.Put("/task/{id}", helpers.MakeHandler(h.HandleUpdateTask))
	r.Post("/task/{id}", helpers.MakeHandler(h.HandleUpdateTask))

	// delete the task
	r.Delete("/task/{id}", helpers.MakeHandler(h.HandleDeleteTask))

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application is running", "port", port)
	return http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), r)
}
