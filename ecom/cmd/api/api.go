package api

import (
	"log"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/dkr290/go-advanced-projects/ecom/handlers/userhandlers"
	"github.com/go-chi/chi/v5"
)

type ApiServer struct {
	addr string
	db   *db.MysqlDB
}

func New(addr string, db *db.MysqlDB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (a *ApiServer) Run() error {
	r := chi.NewRouter()
	h := userhandlers.NewUserHandler(*a.db)
	r.Route("/api/v1", func(r chi.Router) {
		h.RegisterRoutes(r)

	})

	log.Println("Listening on", a.addr)

	return http.ListenAndServe(a.addr, r)
}
