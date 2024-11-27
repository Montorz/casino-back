package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/logger"
	"casino-back/internal/app/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WebSocketHandler struct {
	userService *service.UserService
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketHandler(userService *service.UserService) *WebSocketHandler {
	return &WebSocketHandler{userService: userService}
}

type tokenClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

const signingKey = "secret_key"

func ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.InfoKV("service err", "err", fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		logger.InfoKV("service error", "err", "token claims are not of type *tokenClaims")
		return 0, fmt.Errorf("invalid token or token claims")
	}

	return claims.UserId, nil
}

func (h *WebSocketHandler) StreamBalance(ctx *gin.Context) {
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.InfoKV("error", "err", err)
		}
	}()

	var token string
	_, message, err := conn.ReadMessage()
	if err != nil {
		logger.InfoKV("error", "err", err)
		return
	}

	var response dto.AuthResponse
	if err := json.Unmarshal(message, &response); err != nil || response.Type != "authenticate" {
		logger.InfoKV("error", "err", err)
		return
	}

	token = response.Token
	userId, err := ParseToken(token)
	if err != nil {
		logger.InfoKV("error", "err", err)
	}

	for {
		balance, err := h.userService.GetUserBalance(userId)
		if err != nil {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Error fetching balance")); err != nil {
				logger.InfoKV("service error", "err", err)
				return
			}
			return
		}

		if err := conn.WriteJSON(map[string]interface{}{"balance": balance}); err != nil {
			logger.InfoKV("service error", "err", err)
			return
		}

		time.Sleep(2 * time.Second)
	}
}
