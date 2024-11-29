package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/service"
	"casino-back/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type AuthHandler struct {
	userService     *service.UserService
	jwtTokenManager *token.JwtTokenManager
}

func NewAuthHandler(userService *service.UserService, jwtTokenManager *token.JwtTokenManager) *AuthHandler {
	return &AuthHandler{userService: userService, jwtTokenManager: jwtTokenManager}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var request dto.UserRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.AvatarURL = "/uploads/avatars/default.png"
	userId, err := h.userService.CreateUser(request.Name, request.Login, request.Password, request.AvatarURL)
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

	userId, err := h.userService.GetUserId(response.Login, response.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	claims := &dto.JwtUserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(h.jwtTokenManager.Lifetime))),
		},
		UserId: userId,
	}

	tokenString, err := h.jwtTokenManager.Generate(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User successfully authenticated",
		"token":   tokenString,
	})
}
