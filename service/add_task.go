package service

import (
	"context"
	"fmt"

	"github.com/yoiyoicho/go_todo_app/entity"
	"github.com/yoiyoicho/go_todo_app/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	if err := a.Repo.AddTask(ctx, a.DB, t); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}
