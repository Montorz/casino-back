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

	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, transactionRepository)
	userHandler := handler.NewUserHandler(userService, transactionService)

	authService := service.NewAuthService(userRepository)
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
		api.POST("/transaction", userHandler.ChangeBalance)
	}

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logger.FatalKV("server error", "err", err)
	}
}
