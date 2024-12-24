package handler

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
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
	for {
		// 1. Пауза перед стартом игры
		time.Sleep(10 * time.Second)

		// 2. Генерация результата
		result, err := h.gameService.GetGameResult(gameName)
		if err != nil {
			continue // Пропускаем итерацию при ошибке
		}

		var gameResult dto.GameResultMessage
		var animationDuration time.Duration

		switch gameName {
		case "crash":
			// Получаем результат для игры "Crash"
			crashResult, _ := result.(float64)

			// Вычисляем длительность анимации
			animationDuration = time.Duration(crashResult*float64(time.Second)) + 3*time.Second

			// Формируем сообщение "start"
			gameResult = dto.GameResultMessage{
				Name:   gameName,
				Result: "start",
			}

		case "wheel":
			// Получаем результат для игры "Wheel"
			//wheelResult, _ := result.(int)

			// Генерируем случайную длительность вращения (5-20 секунд)
			randomSeconds := rand.Intn(16) + 5
			animationDuration = time.Duration(randomSeconds)*time.Second + 3*time.Second

			// Формируем сообщение "start"
			gameResult = dto.GameResultMessage{
				Name:   gameName,
				Result: "start",
			}

		default:
			continue
		}

		// 3. Уведомляем клиентов о старте игры
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

		// 4. Ждем окончания "анимации"
		time.Sleep(animationDuration)

		// 5. Отправляем сообщение "stop"
		gameResult.Result = "stop"
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
	}
}
