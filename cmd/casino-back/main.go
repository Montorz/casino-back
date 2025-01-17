package main

import (
	"casino-back/internal/app/handler"
	"casino-back/internal/app/logger"
	"casino-back/internal/app/repository"
	"casino-back/internal/app/service"
	"casino-back/pkg/middleware"
	"casino-back/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=dev dbname=main sslmode=disable")
	if err != nil {
		logger.ErrorKV("db error", "err", err)
	}
	logger.InfoKV("db connected")

	jwtTokenManager := token.NewJwtTokenGenerator("secret_key", 24)

	userRepository := repository.NewUserRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	gameRepository := repository.NewGameRepository(db)

	userService := service.NewUserService(userRepository)
	transactionService := service.NewTransactionService(transactionRepository)
	gameService := service.NewGameService(gameRepository)

	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService, userService)
	gameHandler := handler.NewGameHandler(gameService, userService)

	authHandler := handler.NewAuthHandler(userService, jwtTokenManager)
	webSocketHandler := handler.NewWebSocketHandler(userService, gameService)

	r := gin.New()
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	r.Use(corsConfig)
	r.Static("/uploads", "./uploads")

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}

	ws := r.Group("/ws")
	{
		go webSocketHandler.StartChatBroadcasting()
		ws.GET("/chat", middleware.JwtWebSocketAuth(jwtTokenManager), webSocketHandler.HandleChatWebSocket)
		ws.GET("/:name", webSocketHandler.HandleGameWebSocket)
	}

	api := r.Group("/api", middleware.JwtAuth(jwtTokenManager))
	{
		account := api.Group("/account")
		{
			account.GET("/data", userHandler.GetUserData)
			account.GET("/balance", userHandler.GetUserBalance)
			account.POST("/avatar", userHandler.UpdateAvatar)
		}

		transaction := api.Group("/transaction")
		{
			transaction.POST("/create", transactionHandler.CreateTransaction)
			transaction.GET("/history", transactionHandler.GetTransactions)
		}

		game := api.Group("/game")
		{
			game.GET("/:name/result", gameHandler.GetGameResult)
			game.POST("/bet", gameHandler.PlaceBet)
			game.POST("/create", gameHandler.CreateGame)
			game.GET("/history", gameHandler.GetGames)
		}
	}

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logger.FatalKV("server error", "err", err)
	}
}
