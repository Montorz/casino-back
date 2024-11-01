package main

import (
	"casino-back/internal/app/handler"
	"casino-back/internal/app/logger"
	"casino-back/internal/app/repository"
	"casino-back/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=dev dbname=postgres sslmode=disable")
	if err != nil {
		logger.ErrorKV("db error", "err", err)
	}
	logger.InfoKV("db connected")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	r := gin.New()
	//r.Use(cors.Default())

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", userHandler.CreateUser)
		//auth.POST("/sign-in")
	}

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logger.FatalKV("server error", "err", err)
	}
}
