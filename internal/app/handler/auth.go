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

const defaultAvatarURL = "/uploads/avatars/default.png"

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var request dto.SignUpRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	request.AvatarURL = defaultAvatarURL
	userId, err := h.userService.CreateUser(request.Name, request.Login, request.Password, request.AvatarURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't create user"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User successfully registered",
		"user_id": userId,
	})
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var request dto.SignInRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	userId, err := h.userService.GetUserId(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't get user"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't generate token"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User successfully authenticated",
		"token":   tokenString,
	})
}
