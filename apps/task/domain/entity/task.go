package entity

type Task struct {
	ID     int64
	UserID int64

	Title   string
	Content string
	Status  Status

	CreatedAt int64
	UpdatedAt int64
}

type Status int32

const (
	ToDoStatus Status = iota
	FinishedStatus
)

func (s Status) String() string {
	switch s {
	case ToDoStatus:
		return "todo"
	case FinishedStatus:
		return "finished"
	default:
		return "unknown"
	}
}

func (s Status) Int32() int32 {
	return int32(s)
}
