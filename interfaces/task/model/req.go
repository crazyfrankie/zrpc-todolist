package model

type CreateTaskReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateTaskReq struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}
