package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/handlers"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/middleware"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

var port = ":8080"

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

func main() {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))
	teacherHandler := handlers.NewTeachersHandler(teachers)

	huma.Get(api, "/", teacherHandler.RootHandler)
	huma.Get(api, "/teachers", teacherHandler.TeachersGet)
	huma.Get(api, "/teacher/{id}", teacherHandler.TeacherGet)
	huma.Post(api, "/teachers", teacherHandler.TeachersAdd)
	rl := middleware.NewRateLimit(200, time.Minute)
	server := &http.Server{
		Addr: port,
		Handler: rl.Middleware(middleware.ResponseTimeMiddleware(
			middleware.SecurityHeaders(middleware.Cors(router))),
		),
	}

	fmt.Println("The server is starting on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Error Starting the server", err)
	}
}
