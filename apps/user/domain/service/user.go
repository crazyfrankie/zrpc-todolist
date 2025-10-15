package service

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/entity"
)

type CreateUserRequest struct {
	Name     string
	Password string
}

type User interface {
	Create(ctx context.Context, req *CreateUserRequest) (*entity.User, error)
	Login(ctx context.Context, name, password string) (*entity.User, error)
	ResetPassword(ctx context.Context, name, password string) error
	GetUserInfo(ctx context.Context, userID int64) (user *entity.User, err error)
	UpdateAvatar(ctx context.Context, userID int64, ext string, imagePayload []byte) (url string, err error)
}
