package handler

import (
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SlotHandler struct {
	slotService *service.SlotService
	userService *service.UserService
}

func NewSlotHandler(slotService *service.SlotService, userService *service.UserService) *SlotHandler {
	return &SlotHandler{slotService: slotService, userService: userService}
}

func (h *SlotHandler) GetSlotData(ctx *gin.Context) {
	slotName := ctx.Param("name")
	if slotName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "slot name is required"})
		return
	}

	slotData, err := h.slotService.GetSlotData(slotName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"slotData": slotData,
	})
}

func (h *SlotHandler) PlaceBet(ctx *gin.Context) {
	// Получаем данные о слоте
	slotName := ctx.Param("name")
	if slotName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "slot name is required"})
		return
	}

	slotData, err := h.slotService.GetSlotData(slotName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	minBet := slotData.MinBet
	maxBet := slotData.MaxBet

	// Получаем данные о юзере
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId type"})
		return
	}

	// Обрабатываем запрос ставки
	var input struct {
		BetAmount int `json:"betAmount" binding:"required"`
	}

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.BetAmount < minBet || input.BetAmount > maxBet {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bet amount must be between minBet and maxBet"})
		return
	}

	// Производим вычет ставки из баланса
	err = h.userService.WithdrawBalance(userIDInt, input.BetAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем ответ
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "completed",
	})

}

func (h *SlotHandler) GetSlotResult(ctx *gin.Context) {
	// Получаем данные о слоте
	slotName := ctx.Param("name")
	if slotName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "slot name is required"})
		return
	}

	// В зависимости от слота будем генерировать разные ответы
	switch slotName {
	case "crash":
		crashPoint, err := h.slotService.GetCrashResult()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"crashPoint": crashPoint,
		})
	}
}

func (h *SlotHandler) GetUserResult(ctx *gin.Context) {
	var input struct {
		BetAmount   float64 `json:"betAmount" binding:"required"`
		Coefficient float64 `json:"coefficient" binding:"required"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	winAmount := int(input.BetAmount * input.Coefficient)
	err := h.userService.TopUpBalance(userIDInt, winAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "completed",
		"winAmount": winAmount,
	})
}
