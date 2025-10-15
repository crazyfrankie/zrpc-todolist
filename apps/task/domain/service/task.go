package service

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/entity"
)

type CreateTaskRequest struct {
	UserID  int64
	Title   string
	Content string
}

type UpdateTaskRequest struct {
	TaskID  int64
	Title   *string
	Content *string
}

type Task interface {
	Create(ctx context.Context, req *CreateTaskRequest) (*entity.Task, error)
	GetTaskList(ctx context.Context, userID int64) ([]*entity.Task, error)
	UpdateTask(ctx context.Context, req *UpdateTaskRequest) error
	UpdateTaskStatus(ctx context.Context, taskID int64, status int32) error
	GetTaskRecycleList(ctx context.Context, userID int64) ([]*entity.Task, error)
}
