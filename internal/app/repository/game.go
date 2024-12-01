package repository

import (
	"casino-back/internal/app/logger"
	"casino-back/internal/app/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type GameRepository struct {
	db *sqlx.DB
}

func NewGameRepository(db *sqlx.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) CreateGame(userId int, game model.Game) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id, name, bet_amount, coefficient, win_amount, created_date) values ($1, $2, $3, $4, $5, $6) RETURNING id", "games")
	row := r.db.QueryRow(query, userId, game.Name, game.BetAmount, game.Coefficient, game.WinAmount, game.CreatedDate)

	if err := row.Scan(&id); err != nil {
		logger.InfoKV("game repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *GameRepository) GetGames(userId int) ([]model.Game, error) {
	var game []model.Game

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", "games")
	err := r.db.Select(&game, query, userId)

	if err != nil {
		logger.InfoKV("game repository error", "err", err)
		return nil, err
	}

	return game, nil
}
