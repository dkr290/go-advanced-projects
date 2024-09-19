package userhandlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/config"
	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/dkr290/go-advanced-projects/ecom/helpers"
	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandlers struct {
	DB db.Database
}

func NewUserHandler(db db.Database) *UserHandlers {
	return &UserHandlers{
		DB: db,
	}
}

func (h *UserHandlers) RegisterRoutes(router chi.Router) {
	router.Post("/login", h.handleLogin)
	router.Post("/register", h.handleRegister)
	router.Get("/test", h.testHandler)
}

func (h *UserHandlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := helpers.ParseJson(r, &payload); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
	}

	if err := helpers.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.DB.GetUserByEmail(payload.Email)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	// if it is not nil so passwords dont match returns false !false = true and we handle error
	if !helpers.ComparePasswords(u.Password, []byte(payload.Password)) {
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := helpers.CreateJWT(secret, u.ID)

	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errors.New("create JWT error: "+err.Error()))
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandlers) handleRegister(w http.ResponseWriter, r *http.Request) {
	//get json payload
	//check if the user exists

	//if it does not we create new user
	var payload types.RegisterUserPayload
	if err := helpers.ParseJson(r, &payload); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
	}

	if err := helpers.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	// check if the user exists
	_, err := h.DB.GetUserByEmail(payload.Email)
	if err == nil {
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s alreayd exists", payload.Email))
		return
	}
	hashedPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	//if the user does not exists
	err = h.DB.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := helpers.WriteJSON(w, http.StatusCreated, nil); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *UserHandlers) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Test")
}
