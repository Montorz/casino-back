package handler

import (
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
	var input struct {
		Name        string  `json:"name"`
		BetAmount   float64 `json:"betAmount"`
		Coefficient float64 `json:"coefficient"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем id пользователя
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId type"})
		return
	}

	// Получаем id слота по имени
	slotId, err := h.slotService.GetSlot(input.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Добавляем баланс в зависимости от коэф.
	winAmount := int(input.BetAmount * input.Coefficient)
	err = h.userService.TopUpBalance(userIDInt, winAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	game := model.Game{
		UserId:      userIDInt,
		SlotId:      slotId,
		Name:        input.Name,
		BetAmount:   int(input.BetAmount),
		Coefficient: input.Coefficient,
		WinAmount:   winAmount,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	gameID, err := h.gameService.CreateGame(userIDInt, slotId, game)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "transaction logging failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "completed",
		"game_id": gameID,
	})
}

func (h *GameHandler) GetGames(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId type"})
		return
	}

	gameData, err := h.gameService.GetGames(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"game": gameData,
	})
}
