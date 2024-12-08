package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Room struct {
	results chan dto.GameResultMessage
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

type WebSocketHandler struct {
	userService *service.UserService
	gameService *service.GameService

	rooms       map[string]*Room
	chatClients map[*websocket.Conn]bool
	broadcast   chan dto.ChatMessage
	mu          sync.Mutex
	upgrader    websocket.Upgrader
}

func NewWebSocketHandler(userService *service.UserService, gameService *service.GameService) *WebSocketHandler {
	return &WebSocketHandler{
		userService: userService,
		gameService: gameService,

		rooms:       make(map[string]*Room),
		chatClients: make(map[*websocket.Conn]bool),
		broadcast:   make(chan dto.ChatMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WebSocketHandler) HandleChatWebSocket(ctx *gin.Context) {
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
	h.chatClients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.chatClients, conn)
		h.mu.Unlock()
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			chatMessage := dto.ChatMessage{
				Name:      name,
				AvatarURL: avatarURL,
				Message:   string(msg),
			}
			h.broadcast <- chatMessage
		}
	}
}

func (h *WebSocketHandler) StartChatBroadcasting() {
	for {
		message := <-h.broadcast
		h.mu.Lock()
		for client := range h.chatClients {
			err := client.WriteJSON(message)
			if err != nil {
				err = client.Close()
				if err != nil {
					return
				}
				delete(h.chatClients, client)
			}
		}
		h.mu.Unlock()
	}
}

func (h *WebSocketHandler) HandleGameWebSocket(ctx *gin.Context) {
	gameName := ctx.Param("name")

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to upgrade connection"})
		return
	}

	h.mu.Lock()
	var room *Room
	room, exists := h.rooms[gameName]
	if !exists {
		room = &Room{
			results: make(chan dto.GameResultMessage),
			clients: make(map[*websocket.Conn]bool),
		}
		h.rooms[gameName] = room
		go h.startRoomResults(room, gameName)
	}
	room.mu.Lock()
	room.clients[conn] = true
	room.mu.Unlock()
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(room.clients, conn)
		h.mu.Unlock()
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			room.mu.Lock()
			for client := range room.clients {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					err = client.Close()
					if err != nil {
						return
					}
					delete(room.clients, client)
				}
			}
			room.mu.Unlock()
		}
	}
}

func (h *WebSocketHandler) startRoomResults(room *Room, gameName string) {
	var ticker *time.Ticker
	var tickerInterval time.Duration

	var crashResult float64
	var wheelResult int

	ticker = time.NewTicker(5 * time.Second) // Создаем тикер с начальным интервалом
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := h.gameService.GetGameResult(gameName)
			if err != nil {
				continue
			}

			var gameResult dto.GameResultMessage
			switch gameName {
			case "crash":
				crashResult, _ = result.(float64)

				gameResult = dto.GameResultMessage{
					Name:   gameName,
					Result: crashResult,
				}
			case "wheel":
				wheelResult, _ = result.(int)

				gameResult = dto.GameResultMessage{
					Name:   gameName,
					Result: wheelResult,
				}
			default:
				continue
			}

			room.mu.Lock()
			for client := range room.clients {
				err := client.WriteJSON(gameResult)
				if err != nil {
					err = client.Close()
					if err != nil {
						return
					}
					delete(room.clients, client)
				}
			}
			room.mu.Unlock()

			// Динамически меняем интервал на основе результата
			switch gameName {
			case "crash":
				// 0.1x = 0.1sec
				tickerInterval = time.Duration(crashResult*float64(time.Second)) + 15*time.Second
			case "wheel":
				// 15sec
				tickerInterval = 15 * time.Second
			}

			// Перезапускаем тикер с новым интервалом
			ticker.Stop()
			ticker = time.NewTicker(tickerInterval)
		}
	}
}
