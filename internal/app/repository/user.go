package repository

import (
	"casino-back/internal/app/logger"
	"casino-back/internal/app/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user model.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, login, password, balance) values ($1, $2, $3, $4) RETURNING id", "users")
	row := r.db.QueryRow(query, user.Name, user.Login, user.Password, user.Balance)

	if err := row.Scan(&id); err != nil {
		logger.ErrorKV("repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUser(login string, password string) (int, error) {
	var id int

	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password=$2", "users")
	err := r.db.Get(&id, query, login, password)

	if err != nil {
		logger.ErrorKV("repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetBalance(userId int) (int, error) {
	var balance int

	query := fmt.Sprintf("SELECT balance FROM %s WHERE id=$1", "users")
	err := r.db.Get(&balance, query, userId)

	if err != nil {
		logger.ErrorKV("repository error", "err", err)
		return 0, err
	}

	return balance, nil
}

func (r *UserRepository) UpdateBalance(userId int, newBalance int) error {
	query := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", "users")
	_, err := r.db.Exec(query, newBalance, userId)

	if err != nil {
		logger.ErrorKV("repository error", "err", err)
		return err
	}

	return nil
}
