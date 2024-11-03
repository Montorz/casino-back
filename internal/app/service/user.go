package service

import (
	"casino-back/internal/app/model"
	"fmt"
)

type IUserRepository interface {
	CreateUser(user model.User) (int, error)
	GetUser(login string, password string) (int, error)
	GetBalance(userId int) (int, error)
	UpdateBalance(userId int, newBalance int) error
	GetUserData(userId int) (model.User, error)
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetBalance(userId int) (int, error) {
	return s.userRepository.GetBalance(userId)
}

func (s *UserService) UpdateBalance(userId int, newBalance int) error {
	return s.userRepository.UpdateBalance(userId, newBalance)
}

func (s *UserService) TopUpBalance(userId int, amount int) error {
	balance, err := s.GetBalance(userId)
	if err != nil {
		return err
	}
	newBalance := balance + amount

	return s.userRepository.UpdateBalance(userId, newBalance)
}

func (s *UserService) WithdrawBalance(userId int, amount int) error {
	balance, err := s.GetBalance(userId)
	if err != nil {
		return err
	}

	if balance < amount {
		return fmt.Errorf("insufficient balance: available %d, requested %d", balance, amount)
	}
	newBalance := balance - amount

	return s.userRepository.UpdateBalance(userId, newBalance)
}

func (s *UserService) GetUserData(userId int) (model.User, error) {
	return s.userRepository.GetUserData(userId)
}
