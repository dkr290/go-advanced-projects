package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/frontend/models"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/view/home"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/view/todo"
	"github.com/go-chi/chi"
)

func (h *Handlers) HandleHome(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}

func (h *Handlers) HandleFetchTasks(w http.ResponseWriter, r *http.Request) error {
	// Make a GET request to the database microservice
	resp, err := http.Get("http://backend:8080/tasks")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// check for the response code from the backend
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch tasks, status code: %d", resp.StatusCode)
	}
	var tasks []models.JsonTask
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return err
	}
	return todo.TodoList(tasks).Render(r.Context(), w)
}

func (h *Handlers) HandleGetTaskForm(w http.ResponseWriter, r *http.Request) error {
	return todo.AddTaskForm().Render(r.Context(), w)
}

func (h *Handlers) HandleAddTask(w http.ResponseWriter, r *http.Request) error {
	task := r.FormValue("task")
	// Create a JSON payload
	payload, err := json.Marshal(map[string]string{"task": task})
	if err != nil {
		return err
	}
	// Make a POST request to the database microservice
	resp, err := http.Post(
		"http://backend:8080/addTask",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add task, status code: %d", resp.StatusCode)
	}
	// Make a GET request to the database microservice
	resp, err = http.Get("http://backend:8080/tasks")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// check for the response code from the backend
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch tasks, status code: %d", resp.StatusCode)
	}
	var tasks []models.JsonTask
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return err
	}
	// return a frash list of tasks again to the end user
	return todo.TodoList(tasks).Render(r.Context(), w)
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
