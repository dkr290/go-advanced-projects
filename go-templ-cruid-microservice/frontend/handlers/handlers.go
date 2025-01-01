package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/frontend/models"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/frontend/view/home"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/frontend/view/todo"
	"github.com/go-chi/chi"
)

func (h *Handlers) HandleHome(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}

func (h *Handlers) HandleFetchTasks(w http.ResponseWriter, r *http.Request) error {
	// Make a GET request to the database microservice
	resp, err := http.Get("http://" + h.BackendService + "/tasks")
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
		"http://"+h.BackendService+"/addTask",
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
	resp, err = http.Get("http://" + h.BackendService + "/tasks")
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

	// Make a GET request to the backend microservice to fetch the task
	resp, err := http.Get(fmt.Sprintf("http://"+h.BackendService+"/task/%d", taskID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch task, status code: %d", resp.StatusCode)
	}
	// now we get returned tasks and decode them in json back again
	var task models.JsonTask
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return err
	}

	return todo.UpdateTaskForm(&task).Render(r.Context(), w)
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

	task := models.JsonTask{
		Id:   taskId,
		Task: taskItem,
		Done: taskStatus,
	}
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}
	// Make a PUT request to the backend microservice
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://"+h.BackendService+"/task/%d", taskId),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to update task, status code: %d", resp.StatusCode)
	}
	// Make a GET request to the database microservice
	resp, err = http.Get("http://" + h.BackendService + "/tasks")
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

func (h *Handlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(tid)
	if err != nil {
		return fmt.Errorf("error converting the id %v", err)
	}
	// Make a DELETE request to the backend microservice
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://"+h.BackendService+"/task/%d", taskId),
		nil,
	)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to delete task, status code: %d", resp.StatusCode)
	}
	// Make a GET request to the database microservice
	resp, err = http.Get("http://" + h.BackendService + "/tasks")
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
	// return a fresh list of tasks again to the end user
	return todo.TodoList(tasks).Render(r.Context(), w)
}
