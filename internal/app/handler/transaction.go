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
	transactionService *service.TransactionService
	userService        *service.UserService
}

func NewTransactionHandler(transactionService *service.TransactionService, userService *service.UserService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService, userService: userService}
}

func (h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request dto.TransactionRequest
	if err = ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	switch request.Type {
	case "Пополнение":
		err = h.userService.TopUpBalance(userID, request.Amount)
	case "Снятие":
		err = h.userService.WithdrawBalance(userID, request.Amount)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := model.Transaction{
		UserId:      userID,
		Type:        request.Type,
		Amount:      request.Amount,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	transactionID, err := h.transactionService.CreateTransaction(userID, transaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "create transaction failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "Transaction successfully registered",
		"transaction_id": transactionID,
	})
}

func (h *TransactionHandler) GetTransactions(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionData, err := h.transactionService.GetTransactions(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "get transactions failed"})
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
