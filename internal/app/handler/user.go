package handler

import (
	"casino-back/internal/app/model"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService        *service.UserService
	transactionService *service.TransactionService
}

func NewUserHandler(userService *service.UserService, transactionService *service.TransactionService) *UserHandler {
	return &UserHandler{userService: userService, transactionService: transactionService}
}

func (h *UserHandler) ChangeBalance(ctx *gin.Context) {
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
	case "TopUp":
		err = h.userService.TopUpBalance(userIDInt, input.Amount)
	case "Withdraw":
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
		UserId: userIDInt,
		Type:   input.Type,
		Amount: input.Amount,
	}

	transactionID, err := h.transactionService.CreateTransaction(userIDInt, transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "transaction logging failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":         "balance updated successfully",
		"transaction_id": transactionID,
	})
}
