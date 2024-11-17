package handler

import (
	"casino-back/internal/app/handler/dto"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Slot name is required"})
		return
	}

	slotData, err := h.slotService.GetSlotData(slotName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Slot successfully retrieved",
		"slot_data": slotData,
	})
}

func (h *SlotHandler) PlaceBet(ctx *gin.Context) {
	slotName := ctx.Param("name")
	if slotName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Slot name is required"})
		return
	}

	slotData, err := h.slotService.GetSlotData(slotName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	minBet := slotData.MinBet
	maxBet := slotData.MaxBet

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
		return
	}

	var response dto.BetResponse
	if err = ctx.ShouldBindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if response.BetAmount < minBet || response.BetAmount > maxBet {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bet amount must be between minBet and maxBet"})
		return
	}

	err = h.userService.WithdrawBalance(userIDInt, response.BetAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Bet successfully placed, balance update",
	})
}

func (h *SlotHandler) GetSlotResult(ctx *gin.Context) {
	slotName := ctx.Param("name")
	if slotName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Slot name is required"})
		return
	}

	switch slotName {
	case "crash":
		crashPoint, err := h.slotService.GetCrashResult()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "Crash result successfully retrieved",
			"crash_point": crashPoint,
		})
		return
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot name"})
		return
	}
}
