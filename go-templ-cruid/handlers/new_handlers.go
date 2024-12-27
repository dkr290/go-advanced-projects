package handlers

import (
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/pkg/db"
)

type Handlers struct {
	MYDB db.MysqlDatabase
}

func NewHandlers(db db.MysqlDatabase) *Handlers {
	return &Handlers{
		MYDB: db,
	}
}
