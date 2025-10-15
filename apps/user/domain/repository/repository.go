package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/internal/dal"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/internal/dal/model"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, name string) (*model.User, bool, error)
	UpdatePassword(ctx context.Context, name, password string) error
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateAvatar(ctx context.Context, userID int64, iconURI string) error
	CheckUniqueNameExist(ctx context.Context, uniqueName string) (bool, error)
	CreateUser(ctx context.Context, user *model.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return dal.NewUserDao(db)
}
