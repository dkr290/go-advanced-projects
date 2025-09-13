package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/middleware"
)

var port = ":8080"

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type TeachersInput struct {
	Body struct {
		Name string `json:"name" maxLength:"30" example:"Teacher name"       doc:"Name of the teacher"`
		Age  string `json:"age"                 example:"Age of the teacher" doc:"Age"`
		City string `json:"city" maxLength:"20" example:"City like Sofia"    doc:"City name"`
	}
}
type TeachersOutput struct {
	Body struct {
		Name string `json:"name" maxLength:"30" example:"Teacher name"       doc:"Name of the teacher"`
		Age  string `json:"age"                 example:"Age of the teacher" doc:"Age"`
		City string `json:"city" maxLength:"20" example:"City like Sofia"    doc:"City name"`
	}
}

func rootHandler(ctx context.Context, _ *struct{}) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = "Hello from root Handler"
	return resp, nil
}

func teachersGet(ctx context.Context, _ *struct{}) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = "Hello GET Method on Teachers struct"
	return resp, nil
}

func teachersPost(ctx context.Context, input *TeachersInput) (*TeachersOutput, error) {
	resp := &TeachersOutput{}

	resp.Body.Name = input.Body.Name
	resp.Body.Age = input.Body.Age
	resp.Body.City = input.Body.City

	return resp, nil
}

// func teachersHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
//
// 	case http.MethodGet:
// 		w.Write([]byte("Hello GET Method on Teachers struct"))
// 		fmt.Println("Hello GET Method on Teachers struct")
// 	case http.MethodPost:
// 		w.Write([]byte("Hello Post  Method on Teachers struct"))
// 		fmt.Println("Hello Post  Method on Teachers struct")
// 	case http.MethodPatch:
// 		w.Write([]byte("Hello Patch Method on Teachers struct"))
// 		fmt.Println("Hello Patch Method on Teachers struct")
// 	case http.MethodDelete:
// 		w.Write([]byte("Hello Delete  Method on Teachers struct"))
// 		fmt.Println("Hello Delete  Method on Teachers struct")
// 	case http.MethodPut:
// 		w.Write([]byte("Hello Put Method on Teachers struct"))
// 		fmt.Println("Hello Put  Method on Teachers struct")
//
// 	}
// }
//
// func studentsHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
//
// 	case http.MethodGet:
// 		w.Write([]byte("Hello GET Method on Students struct"))
// 		fmt.Println("Hello GET Method on Students struct")
// 	case http.MethodPost:
// 		w.Write([]byte("Hello Post  Method on Students struct"))
// 		fmt.Println("Hello Post  Method on Students struct")
// 	case http.MethodPatch:
// 		w.Write([]byte("Hello Patch Method on Students struct"))
// 		fmt.Println("Hello Patch Method on Students struct")
// 	case http.MethodDelete:
// 		w.Write([]byte("Hello Delete  Method on Students struct"))
// 		fmt.Println("Hello Delete  Method on Students struct")
// 	case http.MethodPut:
// 		w.Write([]byte("Hello Put Method on Students struct"))
// 		fmt.Println("Hello Put  Method on Students struct")
//
// 	}
// }
//
// func execsHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
//
// 	case http.MethodGet:
// 		w.Write([]byte("Hello GET Method on Execs struct"))
// 		fmt.Println("Hello GET Method on Execs struct")
// 	case http.MethodPost:
// 		w.Write([]byte("Hello Post  Method on Execs struct"))
// 		fmt.Println("Hello Post  Method on Execs struct")
// 	case http.MethodPatch:
// 		w.Write([]byte("Hello Patch Method on Execs struct"))
// 		fmt.Println("Hello Patch Method on Execs struct")
// 	case http.MethodDelete:
// 		w.Write([]byte("Hello Delete  Method on Execs struct"))
// 		fmt.Println("Hello Delete  Method on Execs struct")
// 	case http.MethodPut:
// 		w.Write([]byte("Hello Put Method on Execs struct"))
// 		fmt.Println("Hello Put  Method on Execs struct")
//
// 	}
// }

func main() {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))
	huma.Get(api, "/", rootHandler)
	huma.Get(api, "/teachers", teachersGet)
	huma.Post(api, "/teachers", teachersPost)

	rl := middleware.NewRateLimit(5, time.Minute)
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
