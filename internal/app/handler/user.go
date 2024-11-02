package handler

import (
	"casino-back/internal/app/model"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	var input model.User

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.userService.CreateUser(input)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *UserHandler) SignIn(ctx *gin.Context) {
	var input struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.GenerateToken(input.Login, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *UserHandler) UserIdentity(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	userId, err := h.userService.ParseToken(headerParts[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	ctx.Set("userId", userId)
}
