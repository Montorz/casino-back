package dto

type UserRequest struct {
	Name      string `json:"name" binding:"required"`
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	AvatarURL string `json:"avatar_url"`
}

type UserResponse struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}