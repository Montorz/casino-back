package repository

import (
	"casino-back/internal/app/logger"
	"casino-back/internal/app/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

const transactionsTable = "transactions"

func (r *TransactionRepository) CreateTransaction(userId int, transaction model.Transaction) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id, type, amount, created_date) values ($1, $2, $3, $4) RETURNING id", transactionsTable)
	row := r.db.QueryRow(query, userId, transaction.Type, transaction.Amount, transaction.CreatedDate)

	if err := row.Scan(&id); err != nil {
		logger.ErrorKV("create transaction repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *TransactionRepository) GetTransactions(userId int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", transactionsTable)
	err := r.db.Select(&transactions, query, userId)

	if err != nil {
		logger.ErrorKV("get transactions repository error", "err", err)
		return nil, err
	}

	return transactions, nil
}
