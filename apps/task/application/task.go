package application

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/entity"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/service"
	langslice "github.com/crazyfrankie/zrpc-todolist/pkg/lang/slice"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/ctxutil"
	"github.com/crazyfrankie/zrpc-todolist/protocol/task"
)

type TaskApplicationService struct {
	taskDomain service.Task
	task.UnimplementedTaskServiceServer
}

func NewTaskApplicationService(taskDomain service.Task) *TaskApplicationService {
	return &TaskApplicationService{taskDomain: taskDomain}
}

func (t *TaskApplicationService) AddTask(ctx context.Context, req *task.AddTaskRequest) (*task.AddTaskResponse, error) {
	userID := ctxutil.MustGetUserIDFromCtx(ctx)

	newTask, err := t.taskDomain.Create(ctx, &service.CreateTaskRequest{
		UserID:  userID,
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return &task.AddTaskResponse{
		Data: taskDO2DTO(newTask),
	}, nil
}

func (t *TaskApplicationService) ListTasks(ctx context.Context, req *task.ListTasksRequest) (*task.ListTasksResponse, error) {
	userID := ctxutil.MustGetUserIDFromCtx(ctx)

	tasks, err := t.taskDomain.GetTaskList(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &task.ListTasksResponse{
		Data: langslice.Transform(tasks, func(task *entity.Task) *task.Task {
			return taskDO2DTO(task)
		}),
	}, nil
}

func (t *TaskApplicationService) UpdateTask(ctx context.Context, req *task.UpdateTaskRequest) (*task.UpdateTaskResponse, error) {
	err := t.taskDomain.UpdateTask(ctx, &service.UpdateTaskRequest{
		TaskID:  req.GetTaskID(),
		Content: req.Content,
		Title:   req.Title,
	})
	if err != nil {
		return nil, err
	}

	return &task.UpdateTaskResponse{}, nil
}

func (t *TaskApplicationService) UpdateTaskStatus(ctx context.Context, req *task.UpdateTaskStatusRequest) (*task.UpdateTaskStatusResponse, error) {
	err := t.taskDomain.UpdateTaskStatus(ctx, req.GetTaskID(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	return &task.UpdateTaskStatusResponse{}, nil
}

func (t *TaskApplicationService) RecycleBin(ctx context.Context, req *task.RecycleBinRequest) (*task.RecycleBinResponse, error) {
	userID := ctxutil.MustGetUserIDFromCtx(ctx)

	tasks, err := t.taskDomain.GetTaskRecycleList(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &task.RecycleBinResponse{
		Data: langslice.Transform(tasks, func(task *entity.Task) *task.Task {
			return taskDO2DTO(task)
		}),
	}, nil
}

func taskDO2DTO(taskDo *entity.Task) *task.Task {
	return &task.Task{
		TaskID:    taskDo.ID,
		Title:     taskDo.Title,
		Content:   taskDo.Content,
		Status:    taskDo.Status.String(),
		CreatedAt: taskDo.CreatedAt / 1000,
		UpdatedAt: taskDo.UpdatedAt / 1000,
	}
}
