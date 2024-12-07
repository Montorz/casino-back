package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/model"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GameHandler struct {
	gameService *service.GameService
	userService *service.UserService
}

func NewGameHandler(gameService *service.GameService, userService *service.UserService) *GameHandler {
	return &GameHandler{gameService: gameService, userService: userService}
}

func (h *GameHandler) GetGameResult(ctx *gin.Context) {
	gameName := ctx.Param("name")
	if gameName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "game name is required"})
		return
	}

	switch gameName {
	case "crash":
		crashPoint, err := h.gameService.GetGameResult(gameName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid game result"})
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "Crash result successfully retrieved",
			"crash_point": crashPoint,
		})
		return
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid game name"})
		return
	}
}

func (h *GameHandler) PlaceBet(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response dto.BetResponse
	if err = ctx.ShouldBindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad response"})
		return
	}

	err = h.userService.WithdrawBalance(userID, response.BetAmount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Bet successfully placed, balance update",
	})
}

func (h *GameHandler) CreateGame(ctx *gin.Context) {
	var request dto.GameRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	winAmount := int(request.BetAmount * request.Coefficient)
	err = h.userService.TopUpBalance(userID, winAmount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game := model.Game{
		UserId:      userID,
		Name:        request.Name,
		BetAmount:   int(request.BetAmount),
		Coefficient: request.Coefficient,
		WinAmount:   winAmount,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	gameId, err := h.gameService.CreateGame(userID, game)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "create game failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Game successfully registered",
		"game_id": gameId,
	})
}

func (h *GameHandler) GetGames(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameData, err := h.gameService.GetGames(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "get games failed"})
		return
	}

	var gameHistory []dto.GameResponse
	for _, game := range gameData {
		gameHistory = append(gameHistory, dto.GameResponse{
			Id:          game.Id,
			Name:        game.Name,
			BetAmount:   float64(game.BetAmount),
			Coefficient: game.Coefficient,
			WinAmount:   float64(game.WinAmount),
			CreatedDate: game.CreatedDate,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Games successfully retrieved",
		"game_history": gameHistory,
	})
}
