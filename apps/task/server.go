package task

import (
	"context"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/registry"

	"github.com/crazyfrankie/zrpc-todolist/apps/task/application"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/repository"
	"github.com/crazyfrankie/zrpc-todolist/apps/task/domain/service"
	"github.com/crazyfrankie/zrpc-todolist/protocol/task"
)

func Start(ctx context.Context, client *registry.Registry, srv zrpc.ServiceRegistrar) error {
	basic, err := application.Init(ctx, client)
	if err != nil {
		return err
	}
	taskRepo := repository.NewTaskRepository(basic.DB)
	taskDomain := service.NewTaskDomain(&service.Components{
		TaskRepo: taskRepo,
		IDGen:    basic.IDGen,
	})
	appService := application.NewTaskApplicationService(taskDomain)

	task.RegisterTaskServiceServer(srv, appService)

	return nil
}
