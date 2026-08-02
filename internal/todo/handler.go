package todo

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/db"
)

type Handler struct {
	queries db.Querier
}

func NewHandler(queries db.Querier) *Handler {
	return &Handler{
		queries: queries,
	}
}

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,max=200"`
}

type UpdateTodoRequest struct {
	Title     string `json:"title" binding:"required,max=200"`
	Completed *bool  `json:"completed" binding:"required"`
}

type TodoResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTodoResponse(storedTodo db.Todo) TodoResponse {
	return TodoResponse{
		ID:        storedTodo.ID,
		Title:     storedTodo.Title,
		Completed: storedTodo.Completed,
		CreatedAt: storedTodo.CreatedAt.Time,
		UpdatedAt: storedTodo.UpdatedAt.Time,
	}
}

func parseTodoID(context *gin.Context) (int64, bool) {
	idText := context.Param("id")

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "todo ID must be a number",
		})
		return 0, false
	}

	if id < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "todo ID must be greater than zero",
		})
		return 0, false
	}

	return id, true
}

func (handler *Handler) List(context *gin.Context) {
	storedTodos, err := handler.queries.ListTodos(
		context.Request.Context(),
	)
	if err != nil {
		log.Printf("failed to list todos: %v", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list todos",
		})
		return
	}

	response := make([]TodoResponse, 0, len(storedTodos))

	for _, storedTodo := range storedTodos {
		response = append(response, toTodoResponse(storedTodo))
	}

	context.JSON(http.StatusOK, response)
}

func (handler *Handler) Get(context *gin.Context) {
	id, valid := parseTodoID(context)
	if !valid {
		return
	}

	storedTodo, err := handler.queries.GetTodo(
		context.Request.Context(),
		id,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "todo not found",
		})
		return
	}

	if err != nil {
		log.Printf("failed to get todo %d: %v", id, err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get todo",
		})
		return
	}

	context.JSON(http.StatusOK, toTodoResponse(storedTodo))
}

func (handler *Handler) Create(context *gin.Context) {
	var request CreateTodoRequest

	if err := context.ShouldBindJSON(&request); err != nil {
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

	createdTodo, err := handler.queries.CreateTodo(
		context.Request.Context(),
		title,
	)
	if err != nil {
		log.Printf("failed to create todo: %v", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create todo",
		})
		return
	}

	context.JSON(
		http.StatusCreated,
		toTodoResponse(createdTodo),
	)
}

func (handler *Handler) Update(context *gin.Context) {
	id, valid := parseTodoID(context)
	if !valid {
		return
	}

	var request UpdateTodoRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "title and completed are required; title must not exceed 200 characters",
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

	updatedTodo, err := handler.queries.UpdateTodo(
		context.Request.Context(),
		db.UpdateTodoParams{
			ID:        id,
			Title:     title,
			Completed: *request.Completed,
		},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "todo not found",
		})
		return
	}

	if err != nil {
		log.Printf("failed to update todo %d: %v", id, err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update todo",
		})
		return
	}

	context.JSON(http.StatusOK, toTodoResponse(updatedTodo))
}

func (handler *Handler) Delete(context *gin.Context) {
	id, valid := parseTodoID(context)
	if !valid {
		return
	}

	deletedRows, err := handler.queries.DeleteTodo(
		context.Request.Context(),
		id,
	)
	if err != nil {
		log.Printf("failed to delete todo %d: %v", id, err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete todo",
		})
		return
	}

	if deletedRows == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "todo not found",
		})
		return
	}

	context.Status(http.StatusNoContent)
}
