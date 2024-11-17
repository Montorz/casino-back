package handler

import (
	"casino-back/internal/app/handler/dto"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No userId header"})
		return
	}

	userIDInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId type"})
		return
	}

	var request dto.TransactionRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch request.Type {
	case "Пополнение":
		err = h.userService.TopUpBalance(userIDInt, request.Amount)
	case "Снятие":
		err = h.userService.WithdrawBalance(userIDInt, request.Amount)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction := model.Transaction{
		UserId:      userIDInt,
		Type:        request.Type,
		Amount:      request.Amount,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	transactionID, err := h.transactionService.CreateTransaction(userIDInt, transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction logging failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "Transaction successfully registered",
		"transaction_id": transactionID,
	})
}

func (h *TransactionHandler) GetTransactions(ctx *gin.Context) {
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

	transactionData, err := h.transactionService.GetTransactions(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var transactionHistory []dto.TransactionResponse
	for _, transaction := range transactionData {
		transactionHistory = append(transactionHistory, dto.TransactionResponse{
			Id:          transaction.Id,
			Type:        transaction.Type,
			Amount:      transaction.Amount,
			CreatedDate: transaction.CreatedDate,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":             "Transactions successfully retrieved",
		"transaction_history": transactionHistory,
	})
}
