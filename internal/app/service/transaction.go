package service

import "casino-back/internal/app/model"

type ITransactionRepository interface {
	CreateTransaction(userId int, transaction model.Transaction) (int, error)
	GetTransactions(userId int) ([]model.Transaction, error)
}

type TransactionService struct {
	transactionRepository ITransactionRepository
}

func NewTransactionService(transactionRepository ITransactionRepository) *TransactionService {
	return &TransactionService{transactionRepository: transactionRepository}
}

func (s *TransactionService) CreateTransaction(userId int, transaction model.Transaction) (int, error) {
	return s.transactionRepository.CreateTransaction(userId, transaction)
}

func (s *TransactionService) GetTransactions(userId int) ([]model.Transaction, error) {
	return s.transactionRepository.GetTransactions(userId)
}
