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

func (r *TransactionRepository) CreateTransaction(userId int, transaction model.Transaction) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id, title, amount) values ($1, $2, $3) RETURNING id", "transactions")
	row := r.db.QueryRow(query, userId, transaction.Type, transaction.Amount)

	if err := row.Scan(&id); err != nil {
		logger.ErrorKV("repository error", "err", err)
		return 0, err
	}

	return id, nil
}
