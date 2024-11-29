package dto

import "github.com/golang-jwt/jwt/v5"

type UserRequest struct {
	Name      string `json:"name" binding:"required"`
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type UserResponse struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtUserClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"userId"`
}
