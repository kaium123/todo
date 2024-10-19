// Package repository provides the database operations for the todo endpoint.
package repository

import (
	"context"
	"errors"
	"fmt"
	log "github.com/zuu-development/fullstack-examination-2024/internal/log"
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"gorm.io/gorm"
)

// ITodo Todo is the repository for the todo endpoint.
type ITodo interface {
	Create(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, reqParams *model.DeleteRequest) error
	Update(ctx context.Context, todo *model.Todo) error
	Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error)
	FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error)
}

type InitTodoRepository struct {
	Db  *gorm.DB
	Log *log.Logger
}

type todoReceiver struct {
	log *log.Logger
	db  *gorm.DB
}

// NewTodo returns a new instance of the todo repository.
func NewTodo(initTodoRepository *InitTodoRepository) ITodo {
	return &todoReceiver{
		log: initTodoRepository.Log,
		db:  initTodoRepository.Db,
	}
}

func (td *todoReceiver) Create(ctx context.Context, todo *model.Todo) error {
	if err := td.db.Create(todo).Error; err != nil {
		td.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (td *todoReceiver) Update(ctx context.Context, todo *model.Todo) error {
	if err := td.db.Save(todo).Error; err != nil {
		td.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (td *todoReceiver) Delete(ctx context.Context, reqParams *model.DeleteRequest) error {
	result := td.db.Where("id = ?", reqParams.ID).
		Delete(&model.Todo{})
	if result.RowsAffected == 0 {
		td.log.Error(ctx, model.ErrNotFound.Error())
		return model.ErrNotFound
	}

	if result.Error != nil {
		td.log.Error(ctx, result.Error.Error())
		return result.Error
	}

	td.log.Info(ctx, fmt.Sprintf("Deleted todo with id: %d", reqParams.ID))
	return nil
}

func (td *todoReceiver) Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error) {
	var todo *model.Todo
	err := td.db.Where("id = ?", reqParams.ID).
		Take(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}
		td.log.Error(ctx, err.Error())
		return nil, err
	}

	return todo, nil
}

func (td *todoReceiver) FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error) {
	var todos []*model.Todo

	// Build the base query
	query := td.db.Model(&model.Todo{})

	// Filter by task name using LIKE for substring search (if provided)
	if reqParams.Task != "" {
		query = query.Where("task LIKE ?", "%"+reqParams.Task+"%")
	}

	// Optional filtering by status (if provided)
	if reqParams.Status != "" {
		query = query.Where("status = ?", reqParams.Status)
	}

	// Ordering logic:
	// 1. Incomplete tasks (status != 'done') come first.
	// 2. Sort incomplete tasks by priority (high > medium > low).
	// 3. Sort incomplete tasks by created_at in descending order.
	// 4. "Done" tasks should be sorted by updated_at in ascending order.
	query = query.Order(`
			CASE 
				WHEN status != 'done' THEN 0
				ELSE 1
			END ASC,
		
			CASE 
				WHEN status != 'done' THEN
					CASE
						WHEN priority = 'high' THEN 1
						WHEN priority = 'medium' THEN 2
						WHEN priority = 'low' THEN 3
						ELSE 4
					END
				ELSE NULL
			END ASC,
		
			CASE 
				WHEN status != 'done' THEN created_at
				ELSE NULL
			END DESC,
		
			CASE 
				WHEN status = 'done' THEN updated_at
				ELSE NULL
			END ASC
		`)

	// Execute the query to retrieve sorted tasks
	err := query.Find(&todos).Error
	if err != nil {
		td.log.Error(ctx, err.Error())
		return nil, err
	}

	return todos, nil
}
