package api

import (
	"log"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/dkr290/go-advanced-projects/ecom/handlers/producthandlers"
	"github.com/dkr290/go-advanced-projects/ecom/handlers/userhandlers"
	"github.com/go-chi/chi/v5"
)

type ApiServer struct {
	addr string
	db   *db.MysqlDB
	pdb  *db.ProductMysql
}

func New(addr string, db *db.MysqlDB, pdb *db.ProductMysql) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
		pdb:  pdb,
	}
}

func (a *ApiServer) Run() error {
	r := chi.NewRouter()
	h := userhandlers.NewUserHandler(a.db)
	hp := producthandlers.NewProductHandler(a.pdb)
	r.Route("/api/v1/users", func(r chi.Router) {
		h.RegisterRoutes(r)

	})
	r.Route("/api/v1/products", func(r chi.Router) {
		hp.RegisterRoutes(r)
	})

	log.Println("Listening on", a.addr)

	return http.ListenAndServe(a.addr, r)
}
