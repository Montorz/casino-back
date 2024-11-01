package handler

import (
	"casino-back/internal/app/logger"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, status int, message string) {
	logger.ErrorKV(message)
	ctx.AbortWithStatusJSON(status, gin.H{"error": errorResponse{message}})
}
