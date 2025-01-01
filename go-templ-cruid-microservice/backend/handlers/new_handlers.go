package handlers

import "github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/pkg/db"

type Handlers struct {
	MYDB db.TodoDatabase
}

func NewHandlers(db db.TodoDatabase) *Handlers {
	return &Handlers{
		MYDB: db,
	}
}
