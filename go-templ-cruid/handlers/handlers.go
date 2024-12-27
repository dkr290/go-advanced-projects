package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid/view/home"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid/view/todo"
)

func (h *Handlers) HandleHome(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}

func (h *Handlers) HandleFetchTasks(w http.ResponseWriter, r *http.Request) error {
	todos, err := h.MYDB.GetAllTasks()
	if err != nil {
		return err
	}

	return todo.TodoList(todos).Render(r.Context(), w)
}
