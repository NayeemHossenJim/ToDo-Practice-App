package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,max=200"`
}

var todos = []Todo{
	{
		ID:        1,
		Title:     "Learn Go structs",
		Completed: true,
	},
	{
		ID:        2,
		Title:     "Learn Go slices",
		Completed: false,
	},
	{
		ID:        3,
		Title:     "Build a Todo API",
		Completed: false,
	},
	{
		ID:        4,
		Title:     "Learn Gin",
		Completed: true,
	},
}

func pingHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func listTodosHandler(context *gin.Context) {
	context.JSON(http.StatusOK, todos)
}

func getTodoHandler(context *gin.Context) {
	idText := context.Param("id")

	id, parseError := strconv.ParseInt(idText, 10, 64)
	if parseError != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "todo ID must be a number",
		})
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			context.JSON(http.StatusOK, todo)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "todo not found",
	})
}

func nextTodoID() int64 {
	var highestID int64

	for _, todo := range todos {
		if todo.ID > highestID {
			highestID = todo.ID
		}
	}

	return highestID + 1
}

func createTodoHandler(context *gin.Context) {
	var request CreateTodoRequest

	bindingError := context.ShouldBindJSON(&request)
	if bindingError != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "title is required and must not exceed 200 characters",
		})
		return
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "title cannot be empty",
		})
		return
	}

	newTodo := Todo{
		ID:        nextTodoID(),
		Title:     title,
		Completed: false,
	}

	todos = append(todos, newTodo)

	context.JSON(http.StatusCreated, newTodo)
}

func main() {
	router := gin.Default()

	proxyError := router.SetTrustedProxies(nil)
	if proxyError != nil {
		log.Fatal("failed to configure trusted proxies: ", proxyError)
	}

	router.GET("/ping", pingHandler)
	router.GET("/todos", listTodosHandler)
	router.GET("/todos/:id", getTodoHandler)
	router.POST("/todos", createTodoHandler)

	serverError := router.Run(":8080")
	if serverError != nil {
		log.Fatal("server failed to start: ", serverError)
	}
}
