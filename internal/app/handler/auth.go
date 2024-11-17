package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var request dto.UserRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := h.authService.CreateUser(request.Name, request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User successfully registered",
		"user_id": userId,
	})
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var response dto.UserResponse

	if err := ctx.BindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(response.Login, response.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User successfully authenticated",
		"token":   token,
	})
}

func (h *AuthHandler) UserIdentity(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}

	userId, err := h.authService.ParseToken(headerParts[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}

	ctx.Set("userId", userId)
}
