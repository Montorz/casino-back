package service

import (
	"casino-back/internal/app/logger"
	"casino-back/internal/app/model"
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

func (s *UserService) GetUserBalance(userId int) (int, error) {
	return s.userRepository.GetUserBalance(userId)
}

func (s *UserService) UpdateBalance(userId int, newBalance int) error {
	return s.userRepository.UpdateUserBalance(userId, newBalance)
}

func (s *UserService) TopUpBalance(userId int, amount int) error {
	balance, err := s.GetUserBalance(userId)
	if err != nil {
		return err
	}
	newBalance := balance + amount

	return s.userRepository.UpdateUserBalance(userId, newBalance)
}

func (s *UserService) WithdrawBalance(userId int, amount int) error {
	balance, err := s.GetUserBalance(userId)
	if err != nil {
		return err
	}

	if balance < amount {
		logger.InfoKV("service error", "err", fmt.Sprintf("insufficient balance: available %d, requested %d", balance, amount))
		return err
	}
	newBalance := balance - amount

	return s.userRepository.UpdateUserBalance(userId, newBalance)
}

func (s *UserService) GetUserData(userId int) (*model.User, error) {
	return s.userRepository.GetUserData(userId)
}

func (s *UserService) UpdateAvatarURL(userId int, avatarURL string) error {
	return s.userRepository.UpdateAvatarURL(userId, avatarURL)
}
