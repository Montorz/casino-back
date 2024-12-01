package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const contextUserId = "userId"

func getUserID(ctx *gin.Context) (int, error) {
	userId, exists := ctx.Get(contextUserId)
	if !exists {
		return 0, ctx.AbortWithError(http.StatusUnauthorized, errors.New("userId not found"))
	}

	userIDInt, ok := userId.(int)
	if !ok {
		return 0, ctx.AbortWithError(http.StatusUnauthorized, errors.New("userId not found"))
	}

	return userIDInt, nil
}
