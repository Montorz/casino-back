package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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
		Name:      userData.Name,
		Login:     userData.Login,
		AvatarURL: userData.AvatarURL,
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

func (h *UserHandler) UpdateAvatar(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	avatarDir := fmt.Sprintf("./uploads/avatars/%v", userId)
	if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("avatar%s", ext)
	filePath := filepath.Join(avatarDir, newFileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%v/%s", userId, newFileName)

	if err := h.userService.UpdateAvatarURL(userId.(int), avatarURL); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar URL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Avatar successfully updated",
		"avatar_url": avatarURL,
	})
}
