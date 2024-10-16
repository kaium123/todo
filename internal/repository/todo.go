// Package repository provides the database operations for the todo endpoint.
package repository

import (
	"context"
	"fmt"
	log "github.com/zuu-development/fullstack-examination-2024/internal/log"

	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"gorm.io/gorm"
)

// Todo is the repository for the todo endpoint.
type Todo interface {
	Create(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, reqParams *model.DeleteRequest) error
	Update(ctx context.Context, todo *model.Todo) error
	Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error)
	FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error)
}

type todoReceiver struct {
	log *log.Logger
	db  *gorm.DB
}

// NewTodo returns a new instance of the todo repository.
func NewTodo(db *gorm.DB, log *log.Logger) Todo {
	return &todoReceiver{
		log: log,
		db:  db,
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
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		td.log.Error(ctx, err.Error())
		return nil, err
	}

	return todo, nil
}

func (td *todoReceiver) FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error) {
	var todos []*model.Todo

	query := td.db

	if reqParams.Status != "" {
		query = query.Where("status = ?", reqParams.Status)
	}

	if reqParams.Task != "" {
		query = query.Where("task = ?", reqParams.Task)
	}

	err := query.Find(&todos).Error
	if err != nil {
		td.log.Error(ctx, err.Error())
		return nil, err
	}

	return todos, nil
}
