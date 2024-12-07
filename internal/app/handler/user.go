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
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := h.userService.GetUserData(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "get userData failed"})
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
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.userService.GetUserBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "get userBalance failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "User balance successfully retrieved",
		"user_balance": balance,
	})
}

func (h *UserHandler) UpdateAvatar(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "avatar file is required"})
		return
	}

	avatarDir := fmt.Sprintf("./uploads/avatars/%v", userID)
	if err = os.MkdirAll(avatarDir, os.ModePerm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create directory"})
		return
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("avatar%s", ext)
	filePath := filepath.Join(avatarDir, newFileName)

	if err = ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to save file"})
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%v/%s", userID, newFileName)

	if err = h.userService.UpdateAvatarURL(userID, avatarURL); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to update avatar URL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Avatar successfully updated",
		"avatar_url": avatarURL,
	})
}
