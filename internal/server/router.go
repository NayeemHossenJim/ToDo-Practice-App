package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/todo"
)

func NewRouter(
	todoHandler *todo.Handler,
) (*gin.Engine, error) {
	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf(
			"configure trusted proxies: %w",
			err,
		)
	}

	router.GET("/ping", pingHandler)
	router.GET("/todos", todoHandler.List)
	router.GET("/todos/:id", todoHandler.Get)
	router.POST("/todos", todoHandler.Create)
	router.PUT("/todos/:id", todoHandler.Update)
	router.DELETE("/todos/:id", todoHandler.Delete)

	return router, nil
}

func pingHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
