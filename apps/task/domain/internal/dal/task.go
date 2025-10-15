package dal

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/internal/dal/model"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/internal/dal/query"
)

type TaskDao struct {
	query *query.Query
}

func NewTaskDao(db *gorm.DB) *TaskDao {
	return &TaskDao{query: query.Use(db)}
}

func (t *TaskDao) Create(ctx context.Context, task *model.Task) error {
	return t.query.Task.WithContext(ctx).Create(task)
}

func (t *TaskDao) UpdateTask(ctx context.Context, taskID int64, updates map[string]any) error {
	_, err := t.query.Task.WithContext(ctx).Where(
		t.query.Task.ID.Eq(taskID),
	).Updates(updates)
	return err
}

func (t *TaskDao) UpdateTaskStatus(ctx context.Context, taskID int64, status int32) error {
	_, err := t.query.Task.WithContext(ctx).Where(
		t.query.Task.ID.Eq(taskID),
	).Updates(map[string]any{
		"status":     status,
		"updated_at": time.Now().UnixMilli(),
	})
	return err
}

func (t *TaskDao) GetTasksByID(ctx context.Context, userID int64, status int32) ([]*model.Task, error) {
	return t.query.Task.WithContext(ctx).Where(
		t.query.Task.UserID.Eq(userID),
		t.query.Task.Status.Eq(status),
	).Order(t.query.Task.CreatedAt.Desc()).Find()
}
