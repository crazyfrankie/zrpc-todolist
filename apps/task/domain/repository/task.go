package repository

import (
	"context"

	"gorm.io/gorm"
	
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/internal/dal"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/internal/dal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	UpdateTask(ctx context.Context, taskID int64, updates map[string]any) error
	UpdateTaskStatus(ctx context.Context, taskID int64, status int32) error
	GetTasksByID(ctx context.Context, userID int64, status int32) ([]*model.Task, error)
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return dal.NewTaskDao(db)
}
