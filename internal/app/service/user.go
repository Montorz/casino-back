package service

import (
	"casino-back/internal/app/logger"
	"casino-back/internal/app/model"
	"crypto/sha1"
	"fmt"
)

type IUserRepository interface {
	CreateUser(name, login, password, avatarURL string) (int, error)
	GetUserId(login, password string) (int, error)
	GetUserBalance(userId int) (int, error)
	UpdateUserBalance(userId, newBalance int) error
	GetUserData(userId int) (*model.User, error)
	UpdateAvatarURL(userId int, avatarURL string) error
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(name, login, password, avatarURL string) (int, error) {
	password = s.generatePasswordHash(password)
	return s.userRepository.CreateUser(name, login, password, avatarURL)
}

func (s *UserService) GetUserId(login, password string) (int, error) {
	password = s.generatePasswordHash(password)
	return s.userRepository.GetUserId(login, password)
}

func (s *UserService) GetUserBalance(userId int) (int, error) {
	return s.userRepository.GetUserBalance(userId)
}

func (s *UserService) UpdateBalance(userId int, newBalance int) error {
	return s.userRepository.UpdateUserBalance(userId, newBalance)
}

func (s *UserService) GetUserData(userId int) (*model.User, error) {
	return s.userRepository.GetUserData(userId)
}

func (s *UserService) UpdateAvatarURL(userId int, avatarURL string) error {
	return s.userRepository.UpdateAvatarURL(userId, avatarURL)
}

func (s *UserService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (s *UserService) TopUpBalance(userId int, amount int) error {
	balance, err := s.GetUserBalance(userId)
	if err != nil {
		return fmt.Errorf("can't get balance")
	}
	newBalance := balance + amount

	return s.userRepository.UpdateUserBalance(userId, newBalance)
}

func (s *UserService) WithdrawBalance(userId int, amount int) error {
	balance, err := s.GetUserBalance(userId)
	if err != nil {
		return fmt.Errorf("can't get balance")
	}

	if balance < amount {
		logger.ErrorKV("userBalance service error", "err", "insufficient balance")
		return fmt.Errorf("insufficient balance")
	}

	newBalance := balance - amount
	return s.userRepository.UpdateUserBalance(userId, newBalance)
}
