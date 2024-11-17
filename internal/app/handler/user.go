package handler

import (
	"casino-back/internal/app/handler/dto"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
		return
	}

	userData, err := h.userService.GetUserData(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	publicData := dto.UserDataResponse{
		Name:  userData.Name,
		Login: userData.Login,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "User data successfully retrieved",
		"user_data": publicData,
	})
}

func (h *UserHandler) GetUserBalance(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
		return
	}

	balance, err := h.userService.GetBalance(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "User balance successfully retrieved",
		"user_balance": balance,
	})
}
