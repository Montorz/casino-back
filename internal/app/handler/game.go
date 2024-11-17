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
	slotService *service.SlotService
	userService *service.UserService
}

func NewGameHandler(gameService *service.GameService, slotService *service.SlotService, userService *service.UserService) *GameHandler {
	return &GameHandler{gameService: gameService, slotService: slotService, userService: userService}
}

func (h *GameHandler) CreateGame(ctx *gin.Context) {
	var request dto.GameRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
		return
	}

	slotId, err := h.slotService.GetSlot(request.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	winAmount := int(request.BetAmount * request.Coefficient)
	err = h.userService.TopUpBalance(userIDInt, winAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	game := model.Game{
		UserId:      userIDInt,
		SlotId:      slotId,
		Name:        request.Name,
		BetAmount:   int(request.BetAmount),
		Coefficient: request.Coefficient,
		WinAmount:   winAmount,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	gameId, err := h.gameService.CreateGame(userIDInt, slotId, game)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction logging failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Game successfully registered",
		"game_id": gameId,
	})
}

func (h *GameHandler) GetGames(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
		return
	}

	gameData, err := h.gameService.GetGames(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
