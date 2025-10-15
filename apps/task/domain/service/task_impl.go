package service

import (
	"context"
	"fmt"
	"time"

	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/entity"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/internal/dal/model"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/repository"
	"github.com/crazyfrankie/zrpc-todolist/infra/contract/idgen"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/ptr"
)

type Components struct {
	TaskRepo repository.TaskRepository
	IDGen    idgen.IDGenerator
}

type taskImpl struct {
	*Components
}

func NewTaskDomain(c *Components) Task {
	return &taskImpl{c}
}

func (t *taskImpl) Create(ctx context.Context, req *CreateTaskRequest) (*entity.Task, error) {
	id, err := t.IDGen.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	newTask := &model.Task{
		ID:      id,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
		Status:  entity.ToDoStatus.Int32(),
	}

	err = t.TaskRepo.Create(ctx, newTask)
	if err != nil {
		return nil, err
	}

	return taskPO2DO(newTask), nil
}

func (t *taskImpl) GetTaskList(ctx context.Context, userID int64) ([]*entity.Task, error) {
	taskModels, err := t.TaskRepo.GetTasksByID(ctx, userID, entity.ToDoStatus.Int32())
	if err != nil {
		return nil, err
	}

	tasks := make([]*entity.Task, 0, len(taskModels))
	for _, taskModel := range taskModels {
		tasks = append(tasks, taskPO2DO(taskModel))
	}

	return tasks, nil
}

func (t *taskImpl) UpdateTask(ctx context.Context, req *UpdateTaskRequest) error {
	updates := map[string]any{
		"updated_at": time.Now().UnixMilli(),
	}

	if req.Title != nil {
		updates["title"] = ptr.From(req.Title)
	}
	if req.Content != nil {
		updates["content"] = ptr.From(req.Content)
	}

	return t.TaskRepo.UpdateTask(ctx, req.TaskID, updates)
}

func (t *taskImpl) UpdateTaskStatus(ctx context.Context, taskID int64, status int32) error {
	return t.TaskRepo.UpdateTaskStatus(ctx, taskID, status)
}

func (t *taskImpl) GetTaskRecycleList(ctx context.Context, userID int64) ([]*entity.Task, error) {
	taskModels, err := t.TaskRepo.GetTasksByID(ctx, userID, entity.FinishedStatus.Int32())
	if err != nil {
		return nil, err
	}

	tasks := make([]*entity.Task, 0, len(taskModels))
	for _, taskModel := range taskModels {
		tasks = append(tasks, taskPO2DO(taskModel))
	}

	return tasks, nil
}

func taskPO2DO(taskModel *model.Task) *entity.Task {
	return &entity.Task{
		ID:        taskModel.ID,
		UserID:    taskModel.UserID,
		Title:     taskModel.Title,
		Content:   taskModel.Content,
		Status:    entity.Status(taskModel.Status),
		CreatedAt: taskModel.CreatedAt,
		UpdatedAt: taskModel.UpdatedAt,
	}
}
