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

func (r *UserRepository) CreateUser(name, login, password, avatarURL string) (int, error) {
	var id, balance int

	query := fmt.Sprintf("INSERT INTO %s (name, login, password, balance, avatar_url) values ($1, $2, $3, $4, $5) RETURNING id", "users")
	row := r.db.QueryRow(query, name, login, password, balance, avatarURL)

	if err := row.Scan(&id); err != nil {
		logger.InfoKV("user repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUserId(login, password string) (int, error) {
	var id int

	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password=$2", "users")
	err := r.db.Get(&id, query, login, password)

	if err != nil {
		logger.InfoKV("user repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUserBalance(userId int) (int, error) {
	var balance int

	query := fmt.Sprintf("SELECT balance FROM %s WHERE id=$1", "users")
	err := r.db.Get(&balance, query, userId)

	if err != nil {
		logger.InfoKV("user repository error", "err", err)
		return 0, err
	}

	return balance, nil
}

func (r *UserRepository) UpdateUserBalance(userId, newBalance int) error {
	query := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", "users")
	_, err := r.db.Exec(query, newBalance, userId)

	if err != nil {
		logger.InfoKV("user repository error", "err", err)
		return err
	}

	return nil
}

func (r *UserRepository) GetUserData(userId int) (*model.User, error) {
	var user model.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", "users")
	err := r.db.Get(&user, query, userId)

	if err != nil {
		logger.InfoKV("user repository error", "err", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateAvatarURL(userId int, avatarURL string) error {
	query := fmt.Sprintf("UPDATE %s SET avatar_url = $1 WHERE id = $2", "users")
	_, err := r.db.Exec(query, avatarURL, userId)
	if err != nil {
		logger.InfoKV("user repository error", "err", err)
		return err
	}

	return nil
}
