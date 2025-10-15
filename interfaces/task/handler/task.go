package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/interfaces/task/model"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/response"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/conv"
	"github.com/crazyfrankie/zrpc-todolist/protocol/task"
)

type TaskHandler struct {
	taskClient task.TaskServiceClient
}

func NewTaskHandler(taskClient task.TaskServiceClient) *TaskHandler {
	return &TaskHandler{taskClient: taskClient}
}

func (t *TaskHandler) RegisterRoute(r *gin.RouterGroup) {
	taskGroup := r.Group("tasks")
	{
		taskGroup.POST("create", t.CreateTask())
		taskGroup.GET("list", t.ListTask())
		taskGroup.GET("recycle-list", t.RecycleListTask())
		taskGroup.PUT("update/:id", t.UpdateTask())
		taskGroup.PUT("update/:id/status", t.UpdateTaskStatus())
	}
}

func (t *TaskHandler) CreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreateTaskReq
		if err := c.ShouldBind(&req); err != nil {
			response.InvalidParamError(c, err.Error())
			return
		}

		res, err := t.taskClient.AddTask(c.Request.Context(), &task.AddTaskRequest{
			Title:   req.Title,
			Content: req.Content,
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, res.GetData())
	}
}

func (t *TaskHandler) ListTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := t.taskClient.ListTasks(c.Request.Context(), &task.ListTasksRequest{})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, res.GetData())
	}
}

func (t *TaskHandler) RecycleListTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := t.taskClient.RecycleBin(c.Request.Context(), &task.RecycleBinRequest{})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, res.GetData())
	}
}

func (t *TaskHandler) UpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.UpdateTaskReq
		if err := c.ShouldBind(&req); err != nil {
			response.InvalidParamError(c, err.Error())
			return
		}

		taskID, _ := conv.StrToInt64(c.Param("id"))

		_, err := t.taskClient.UpdateTask(c.Request.Context(), &task.UpdateTaskRequest{
			TaskID:  taskID,
			Title:   req.Title,
			Content: req.Content,
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, nil)
	}
}

func (t *TaskHandler) UpdateTaskStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, _ := conv.StrToInt64(c.Param("id"))
		status, _ := conv.StrToInt64(c.Query("status"))

		_, err := t.taskClient.UpdateTaskStatus(c.Request.Context(), &task.UpdateTaskStatusRequest{
			TaskID: taskID,
			Status: int32(status),
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, nil)
	}
}
