package service

import (
	"casino-back/internal/app/model"
	"crypto/sha1"
	"fmt"
)

type IUserRepository interface {
	CreateUser(user model.User) (int, error)
}

type UserService struct {
	userRepository IUserRepository
}

func (s *UserService) CreateUser(user model.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.userRepository.CreateUser(user)
}

func (s *UserService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
