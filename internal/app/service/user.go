package service

import (
	"casino-back/internal/app/model"
)

type IUserRepository interface {
	CreateUser(user model.User) (int, error)
	GetUser(login string, password string) (int, error)
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
