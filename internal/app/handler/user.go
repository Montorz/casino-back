package handler

import (
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

func (h *UserHandler) GetUserData(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId type"})
		return
	}

	userData, err := h.userService.GetUserData(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type PublicUserData struct {
		Name  string `json:"name"`
		Login string `json:"login"`
	}

	publicData := PublicUserData{
		Name:  userData.Name,
		Login: userData.Login,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": publicData,
	})
}

func (h *UserHandler) GetUserBalance(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId type"})
		return
	}

	balance, err := h.userService.GetBalance(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}
