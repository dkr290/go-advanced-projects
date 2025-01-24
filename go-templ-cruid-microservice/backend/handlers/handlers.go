package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/models"
	"github.com/go-chi/chi"
)

// handle calling the DB function to fetch all tasks
// it calls the same DB method but this time encode result to json back to be consumed by frontend
func (h *Handlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	todos, err := h.MYDB.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *Handlers) HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.MYDB.AddTask(task.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	tid := chi.URLParam(r, "id")

	taskId, err := strconv.Atoi(tid)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	task := models.Task{
		Id: taskId,
	}

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// do the update task by ID and also passing the task
	err = h.MYDB.UpdateTaskByID(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	tid := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(tid)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	err = h.MYDB.DeleteTaskByID(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) HandleGetTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.MYDB.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with the task as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) FavIconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func (h *Handlers) TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Invalid request method", r.Method)
		http.Error(w, "Invalid task ID", http.StatusMethodNotAllowed)
		return
	}
	s := "Liveness and Readiness"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(
		w,
		"<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Test for %s</h1>",
		s,
	)
	fmt.Fprintf(w, "</body></html>")
}
