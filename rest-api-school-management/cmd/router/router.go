package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/handlers"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

var (
	teachers = make(map[int]models.Teacher)
	nextID   = 1
)

// initialize dummy data
func init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Algebra",
	}
}

func Router() *http.ServeMux {
	router := http.NewServeMux()
	teacherHandler := handlers.NewTeachersHandler(teachers)

	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Get(api, "/", teacherHandler.RootHandler)
	huma.Get(api, "/teachers", teacherHandler.TeachersGet)
	huma.Get(api, "/teacher/{id}", teacherHandler.TeacherGet)
	huma.Post(api, "/teachers", teacherHandler.TeachersAdd)

	return router
}
