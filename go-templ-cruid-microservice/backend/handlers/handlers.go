package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/models"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/view/todo"
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
	var task models.JsonTask
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

func (h *Handlers) HandleGetTaskUpdateForm(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	taskID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("error converting the id %v", err)
	}

	task, err := h.MYDB.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	return todo.UpdateTaskForm(task).Render(r.Context(), w)
}

func (h *Handlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	taskItem := r.FormValue("task")
	isDone := r.FormValue("done")
	tid := chi.URLParam(r, "id")
	var taskStatus bool

	switch strings.ToLower(isDone) {
	case "yes", "on":
		taskStatus = true
	case "no", "off":
		taskStatus = false
	default:
		taskStatus = false
	}
	taskId, err := strconv.Atoi(tid)
	if err != nil {
		return fmt.Errorf("error converting the id %v, for %v", err, taskItem)
	}

	task := models.Task{
		Id:   taskId,
		Task: taskItem,
		Done: taskStatus,
	}

	// do the update task by ID and also passing the task
	err = h.MYDB.UpdateTaskByID(task)
	if err != nil {
		return err
	}
	todos, err := h.MYDB.GetAllTasks()
	if err != nil {
		return err
	}
	// return a fresh list of tasks again to the end user
	return todo.TodoList(todos).Render(r.Context(), w)
}

func (h *Handlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(tid)
	if err != nil {
		return fmt.Errorf("error converting the id %v", err)
	}
	err = h.MYDB.DeleteTaskByID(taskId)
	if err != nil {
		return err
	}
	todos, err := h.MYDB.GetAllTasks()
	if err != nil {
		return err
	}
	// return a fresh list of tasks again to the end user
	return todo.TodoList(todos).Render(r.Context(), w)
}
