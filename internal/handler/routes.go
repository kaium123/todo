package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/zuu-development/fullstack-examination-2024/internal/log"
	"github.com/zuu-development/fullstack-examination-2024/internal/repository"
	"github.com/zuu-development/fullstack-examination-2024/internal/service"

	"gorm.io/gorm"
)

type ServiceRegistry struct {
	EchoEngine  *echo.Echo
	RedisClient *redis.Client
	DBInstance  *gorm.DB
	Log         *log.Logger
}

// Register registers the routes for the application.
func Register(serviceRegistry *ServiceRegistry) {
	serviceRegistry.EchoEngine.Validator = &CustomValidator{validator: validator.New()}

	api := serviceRegistry.EchoEngine.Group("/api/v1")

	// Health check
	healthHandler := NewHealth()
	api.GET("/healthz", healthHandler.Healthz)

	// Inject Todo Dependency
	redisRepository := repository.NewRedisCache(&repository.InitRedisCache{
		Client: serviceRegistry.RedisClient, Log: serviceRegistry.Log,
	})
	repository := repository.NewTodo(&repository.InitTodoRepository{
		Db: serviceRegistry.DBInstance, Log: serviceRegistry.Log,
	})
	service := service.NewTodo(&service.InitTodoService{
		Log: serviceRegistry.Log, TodoRepository: repository, RedisCache: redisRepository,
	})
	todoHandler := NewTodo(&InitTodoHandler{
		Service: service, Log: serviceRegistry.Log,
	})

	// Add routes for todo
	todo := api.Group("/todos")
	{
		todo.POST("", todoHandler.Create)
		todo.GET("", todoHandler.FindAll)
		todo.GET("/:id", todoHandler.Find)
		todo.PUT("/:id", todoHandler.Update)
		todo.DELETE("/:id", todoHandler.Delete)
	}
}
