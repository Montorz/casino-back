package repository

import (
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
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUser(login string, password string) (int, error) {
	var id int

	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password=$2", "users")
	err := r.db.Get(&id, query, login, password)

	if err != nil {
		return 0, err
	}

	return id, nil
}
