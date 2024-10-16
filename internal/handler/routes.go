package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/zuu-development/fullstack-examination-2024/internal/log"
	"github.com/zuu-development/fullstack-examination-2024/internal/repository"
	"github.com/zuu-development/fullstack-examination-2024/internal/service"
	"gorm.io/gorm"
)

// Register registers the routes for the application.
func Register(e *echo.Echo, db *gorm.DB) {
	e.Validator = &CustomValidator{validator: validator.New()}

	api := e.Group("/api/v1")

	// Health check
	healthHandler := NewHealth()
	api.GET("/healthz", healthHandler.Healthz)

	// Todo
	logger := log.New()
	repository := repository.NewTodo(db, logger)
	service := service.NewTodo(repository, logger)
	todoHandler := NewTodo(service, logger)
	todo := api.Group("/todos")
	{
		todo.POST("", todoHandler.Create)
		todo.GET("", todoHandler.FindAll)
		todo.GET("/:id", todoHandler.Find)
		todo.PUT("/:id", todoHandler.Update)
		todo.DELETE("/:id", todoHandler.Delete)
	}
}
