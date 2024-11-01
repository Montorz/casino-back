package handler

import (
	"casino-back/internal/app/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("invalid json create user"))
		return
	}

	err = h.userService.CreateUser(request.Name, request.Login, request.Password, request.Balance)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("user can't create"))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("user created"))
}
