// Package service provides the business logic for the todo endpoint.
package service

import (
	"context"
	"fmt"
	"github.com/zuu-development/fullstack-examination-2024/internal/log"
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"github.com/zuu-development/fullstack-examination-2024/internal/repository"
)

// Todo is the service for the todo endpoint.
type Todo interface {
	Create(ctx context.Context, reqTodo *model.CreateRequest) (*model.Todo, error)
	Update(ctx context.Context, reqTodo *model.UpdateRequest) (*model.Todo, error)
	Delete(ctx context.Context, reqParams *model.DeleteRequest) error
	Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error)
	FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error)
}

type todoReceiver struct {
	log            *log.Logger
	todoRepository repository.Todo
}

// NewTodo creates a new Todo service.
func NewTodo(r repository.Todo, log *log.Logger) Todo {
	return &todoReceiver{
		log:            log,
		todoRepository: r,
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

	t.log.Info(ctx, fmt.Sprintf("Todo updated successfully with ID: %d", updatedTodo.ID))
	return updatedTodo, nil
}

func (t *todoReceiver) Delete(ctx context.Context, reqParams *model.DeleteRequest) error {
	if err := t.todoRepository.Delete(ctx, reqParams); err != nil {
		t.log.Error(ctx, err.Error())
		return err
	}
	return nil
}

func (t *todoReceiver) Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error) {
	todo, err := t.todoRepository.Find(ctx, reqParams)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return nil, err
	}
	return todo, nil
}

func (t *todoReceiver) FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error) {
	todo, err := t.todoRepository.FindAll(ctx, reqParams)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return nil, err
	}
	return todo, nil
}
