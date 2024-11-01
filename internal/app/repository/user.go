package repository

import (
	"casino-back/internal/app/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(name string, login string, password string, balance int) error {
	user := []model.User{
		{Name: name, Login: login, Password: password, Balance: balance},
	}

	_, err := r.db.NamedExec("INSERT INTO users (name, login, password_hash, balance) VALUES (:name, :login, :password, :balance)", user)
	if err != nil {
		return err
	}

	return nil
}
