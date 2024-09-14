package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/services/user"
	"github.com/go-chi/chi/v5"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func New(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (a *ApiServer) Run() error {
	r := chi.NewRouter()
	userStore := user.NewStore(a.db)
	userHandler := user.NewHandler(&userStore)
	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)

	})

	log.Println("Listening on", a.addr)

	return http.ListenAndServe(a.addr, r)
}
