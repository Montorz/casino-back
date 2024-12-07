package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WebSocketHandler struct {
	userService *service.UserService
	clients     map[*websocket.Conn]bool
	broadcast   chan dto.ChatMessage
	mu          sync.Mutex
	upgrader    websocket.Upgrader
}

func NewWebSocketHandler(userService *service.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		userService: userService,
		clients:     make(map[*websocket.Conn]bool),
		broadcast:   make(chan dto.ChatMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := h.userService.GetUserData(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get user data"})
		return
	}

	name := userData.Name
	avatarURL := userData.AvatarURL

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to upgrade connection"})
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		chatMessage := dto.ChatMessage{
			Name:      name,
			AvatarURL: avatarURL,
			Message:   string(msg),
		}

		h.broadcast <- chatMessage
	}
}

func (h *WebSocketHandler) StartBroadcasting() {
	for {
		message := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteJSON(message)
			if err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}
