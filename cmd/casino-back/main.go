package main

import (
	"casino-back/internal/app/handler"
	"casino-back/internal/app/logger"
	"casino-back/internal/app/repository"
	"casino-back/internal/app/service"
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

	userRepository := repository.NewUserRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	userService := service.NewUserService(userRepository)
	transactionService := service.NewTransactionService(transactionRepository)
	authService := service.NewAuthService(userRepository)

	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(userService, transactionService)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.New()
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	r.Use(corsConfig)

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}

	api := r.Group("/api", authHandler.UserIdentity)
	{
		account := api.Group("/account")
		{
			account.POST("", userHandler.GetUserData)
			account.POST("/balance", userHandler.GetUserBalance)
		}

		transaction := api.Group("/transaction")
		{
			transaction.POST("/create", transactionHandler.CreateTransaction)
			transaction.POST("/history", transactionHandler.GetTransactions)
		}
	}

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logger.FatalKV("server error", "err", err)
	}
}
