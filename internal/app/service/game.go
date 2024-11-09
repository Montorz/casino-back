package service

import "casino-back/internal/app/model"

type IGameRepository interface {
	CreateGame(userId int, slotId int, game model.Game) (int, error)
	GetGames(userId int) ([]model.Game, error)
}

type GameService struct {
	gameRepository IGameRepository
}

func NewGameService(gameRepository IGameRepository) *GameService {
	return &GameService{gameRepository: gameRepository}
}

func (s *GameService) CreateGame(userId int, slotId int, game model.Game) (int, error) {
	return s.gameRepository.CreateGame(userId, slotId, game)
}

func (s *GameService) GetGames(userId int) ([]model.Game, error) {
	return s.gameRepository.GetGames(userId)
}
