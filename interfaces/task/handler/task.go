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

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with title and content
// @Tags Task
// @Accept json
// @Produce json
// @Param request body model.CreateTaskReq true "Create task request"
// @Success 200 {object} response.Response "Task created successfully"
// @Failure 400 {object} response.Response "Invalid parameters"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /tasks/create [post]
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

// ListTask godoc
// @Summary Get task list
// @Description Get all tasks for current user
// @Tags Task
// @Produce json
// @Success 200 {object} response.Response "Task list retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /tasks/list [get]
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

// RecycleListTask godoc
// @Summary Get recycle bin task list
// @Description Get all deleted tasks in recycle bin
// @Tags Task
// @Produce json
// @Success 200 {object} response.Response "Recycle bin task list retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /tasks/recycle-list [get]
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

// UpdateTask godoc
// @Summary Update task
// @Description Update task title and content by task ID
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body model.UpdateTaskReq true "Update task request"
// @Success 200 {object} response.Response "Task updated successfully"
// @Failure 400 {object} response.Response "Invalid parameters"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /tasks/update/{id} [put]
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

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Update task status by task ID
// @Tags Task
// @Produce json
// @Param id path string true "Task ID"
// @Param status query int true "Task status"
// @Success 200 {object} response.Response "Task status updated successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /tasks/update/{id}/status [put]
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
