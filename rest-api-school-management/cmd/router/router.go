package router

import (
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/handlers"
)

func Router(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()
	teachersDB := dataops.NewTeachersDB(db)
	teacherHandler := handlers.NewTeachersHandler(teachersDB)

	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Get(api, "/", teacherHandler.RootHandler)
	huma.Get(api, "/teacher/{id}", teacherHandler.TeacherGet)
	huma.Post(api, "/teachers", teacherHandler.TeachersAdd)
	huma.Get(api, "/teachers", teacherHandler.TeachersGet)

	return router
}
