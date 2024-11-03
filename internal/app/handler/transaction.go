package handler

import (
	"casino-back/internal/app/model"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TransactionHandler struct {
	userService        *service.UserService
	transactionService *service.TransactionService
}

func NewTransactionHandler(userService *service.UserService, transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{userService: userService, transactionService: transactionService}
}

func (h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
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

	var input struct {
		Type   string `json:"type" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch input.Type {
	case "Пополнение":
		err = h.userService.TopUpBalance(userIDInt, input.Amount)
	case "Снятие":
		err = h.userService.WithdrawBalance(userIDInt, input.Amount)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction := model.Transaction{
		UserId:      userIDInt,
		Type:        input.Type,
		Amount:      input.Amount,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	transactionID, err := h.transactionService.CreateTransaction(userIDInt, transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "transaction logging failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":         "completed",
		"transaction_id": transactionID,
	})
}

func (h *TransactionHandler) GetTransactions(ctx *gin.Context) {
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

	transactionData, err := h.transactionService.GetTransactions(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transaction": transactionData,
	})
}
