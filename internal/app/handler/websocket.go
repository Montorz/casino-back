package handler

import (
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WebSocketHandler struct {
	UserService *service.UserService
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketHandler(userService *service.UserService) *WebSocketHandler {
	return &WebSocketHandler{UserService: userService}
}

func (h *WebSocketHandler) StreamBalance(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	userID, exists := c.Get("userId")
	if !exists {
		err := conn.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	for {
		balance, err := h.UserService.GetUserBalance(userID.(int))
		if err != nil {
			err := conn.WriteMessage(websocket.TextMessage, []byte("Error fetching balance"))
			if err != nil {
				return
			}
			return
		}

		err = conn.WriteJSON(map[string]interface{}{"balance": balance})
		if err != nil {
			return
		}

		// Обновляем баланс каждую секунду
		time.Sleep(1 * time.Second)
	}
}
