// Package service provides the business logic for the todo endpoint.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zuu-development/fullstack-examination-2024/internal/log"
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"github.com/zuu-development/fullstack-examination-2024/internal/repository"
)

// Todo is the service for the todo endpoint.
type ITodo interface {
	Create(ctx context.Context, reqTodo *model.CreateRequest) (*model.Todo, error)
	Update(ctx context.Context, reqTodo *model.UpdateRequest) (*model.Todo, error)
	Delete(ctx context.Context, reqParams *model.DeleteRequest) error
	Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error)
	FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error)
}

type todoReceiver struct {
	log            *log.Logger
	todoRepository repository.ITodo
	redisCache     repository.IRedisCache
}

type InitTodoService struct {
	Log            *log.Logger
	TodoRepository repository.ITodo
	RedisCache     repository.IRedisCache
}

// NewTodo creates a new Todo service.
func NewTodo(initTodoService *InitTodoService) ITodo {
	return &todoReceiver{
		log:            initTodoService.Log,
		todoRepository: initTodoService.TodoRepository,
		redisCache:     initTodoService.RedisCache,
	}
}

func (t *todoReceiver) Create(ctx context.Context, reqTodo *model.CreateRequest) (*model.Todo, error) {

	// Create a new Todo instance using the struct-based constructor
	todoModel := model.NewTodo(reqTodo)

	// Validate the input before proceeding
	if err := todoModel.ValidateCreateRequest(); err != nil {
		t.log.Error(ctx, fmt.Sprintf("invalid request: %s", err.Error()))
		return nil, err
	}

	// Attempt to store the new todo using the repository pattern
	if err := t.todoRepository.Create(ctx, todoModel); err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to create todo: %s", err.Error()))
		return nil, err
	}

	todoKey := fmt.Sprintf("todo:%d", todoModel.ID)
	err := t.redisCache.Add(ctx, todoKey, todoModel)
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to create todo: %s", err.Error()))
		return nil, err
	}

	err = t.redisCache.Add(ctx, todoKey, todoModel)
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to create todo: %s", err.Error()))
		return nil, err
	}

	t.log.Info(ctx, fmt.Sprintf("Todo created successfully with ID: %d", todoModel.ID))
	return todoModel, nil
}

func (t *todoReceiver) Update(ctx context.Context, reqTodo *model.UpdateRequest) (*model.Todo, error) {
	// 現在の値を取得
	currentTodo, err := t.Find(ctx, &model.FindRequest{
		ID: reqTodo.ID,
	})
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to find todo with ID: %d and Error: %s", reqTodo.ID, err.Error()))
		return nil, err
	}

	// Update fields only if they are provided in the request
	updatedTodo := model.NewUpdateTodo(reqTodo)
	updatedTodo.PrepareUpdatedTodo(currentTodo)

	// Save updated todo in the repository
	if err := t.todoRepository.Update(ctx, updatedTodo); err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to update  todo with ID: %d and Error: %s", reqTodo.ID, err.Error()))
		return nil, err
	}

	todoKey := fmt.Sprintf("todo:%d", updatedTodo.ID)
	err = t.redisCache.Delete(ctx, todoKey)
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to delete todo in redis with ID: %d and Error: %s", reqTodo.ID, err.Error()))
	}

	err = t.redisCache.Add(ctx, todoKey, updatedTodo)
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("failed to add todo in redis : %s", err.Error()))
	}

	t.log.Info(ctx, fmt.Sprintf("Todo updated successfully with ID: %d", updatedTodo.ID))
	return updatedTodo, nil
}

func (t *todoReceiver) Delete(ctx context.Context, reqParams *model.DeleteRequest) error {
	cacheKey := fmt.Sprintf("todo:%d", reqParams.ID)

	if err := t.todoRepository.Delete(ctx, reqParams); err != nil {
		t.log.Error(ctx, err.Error())
		return err
	}
	err := t.redisCache.Delete(ctx, cacheKey)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return err
	}
	return nil
}
func (t *todoReceiver) Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error) {
	// Try to fetch from Redis cache first
	cacheKey := fmt.Sprintf("todo:%d", reqParams.ID)
	cachedTodo, err := t.redisCache.Get(ctx, cacheKey)
	if err == nil && cachedTodo != "" {
		// If found in Redis, return it
		todo := &model.Todo{}
		err = json.Unmarshal([]byte(cachedTodo), todo)
		if err != nil {
			t.log.Error(ctx, fmt.Sprintf("failed to unmarshal cached todo: %s ", err.Error()))
		} else {
			return todo, nil
		}
	} else if !errors.Is(err, redis.Nil) && err != nil {
		t.log.Error(ctx, fmt.Sprintf("Redis get error: %s", err.Error()))
	}

	// Fetch from the database if not in Redis or unmarshalling failed
	todo, err := t.todoRepository.Find(ctx, reqParams)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return nil, err
	}

	err = t.redisCache.Add(ctx, cacheKey, todo)
	if err != nil {
		t.log.Error(ctx, err.Error())
	}

	return todo, nil
}
func (t *todoReceiver) FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error) {
	//// Check if the results are already cached in Redis
	t.redisCache.DeleteAll(ctx)
	all, err := t.redisCache.FindAll(ctx, reqParams)
	if err == nil {
		return all, nil
	}

	// Cache miss, fetch from the database
	todos, err := t.todoRepository.FindAll(ctx, reqParams)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return nil, err
	}

	// Cache the result in Redis for future requests
	for _, todo := range todos {
		todoKey := fmt.Sprintf("todo:%d", todo.ID)
		err := t.redisCache.Add(ctx, todoKey, todo)
		if err != nil {
			t.log.Error(ctx, err.Error())
			return nil, err
		}
	}

	return todos, nil
}

//func CalculateScore(todo *model.Todo) float64 {
//	// Calculate score based on status, priority, and timestamps
//	var score float64
//
//	// Step 1: Lower score for 'done' status
//	if todo.Status == "done" {
//		score = 0 // Set score to the lowest for done tasks
//	} else {
//		// Step 2: Use created_at timestamp for sorting non-done tasks
//		score += float64(todo.CreatedAt.Unix()) / 100000
//
//		// Step 3: Adjust score based on priority
//		switch todo.Priority {
//		case "high":
//			score += 3.0 // High priority gets the highest score
//		case "medium":
//			score += 2.0 // Medium priority gets a medium score
//		case "low":
//			score += 1.0 // Low priority gets the lowest score
//		}
//	}
//
//	return score
//}
