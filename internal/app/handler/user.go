package handler

import (
	"casino-back/internal/app/model"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var input model.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.userService.CreateUser(input)

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
