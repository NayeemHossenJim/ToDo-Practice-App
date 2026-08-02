package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/database"
	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/db"
	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/todo"
)

func pingHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(
			"could not load .env; using system environment variables",
		)
	}

	databaseURL := os.Getenv("DATABASE_URL")

	startupContext, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := database.NewPostgresPool(
		startupContext,
		databaseURL,
	)
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL: ", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	queries := db.New(pool)
	todoHandler := todo.NewHandler(queries)

	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal("failed to configure trusted proxies: ", err)
	}

	router.GET("/ping", pingHandler)
	router.GET("/todos", todoHandler.List)
	router.GET("/todos/:id", todoHandler.Get)
	router.POST("/todos", todoHandler.Create)
	router.PUT("/todos/:id", todoHandler.Update)
	router.DELETE("/todos/:id", todoHandler.Delete)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("server failed to start: ", err)
	}
}
