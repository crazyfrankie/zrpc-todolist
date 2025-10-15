package dal

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/internal/dal/model"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/internal/dal/query"
)

type UserDao struct {
	query *query.Query
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{query: query.Use(db)}
}

func (u *UserDao) GetUserByName(ctx context.Context, name string) (*model.User, bool, error) {
	user, err := u.query.User.WithContext(ctx).Where(u.query.User.Name.Eq(name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return user, true, err
}

func (u *UserDao) UpdatePassword(ctx context.Context, name, password string) error {
	_, err := u.query.User.WithContext(ctx).Where(
		u.query.User.Name.Eq(name),
	).Updates(map[string]interface{}{
		"password":   password,
		"updated_at": time.Now().UnixMilli(),
	})
	return err
}

func (u *UserDao) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	return u.query.User.WithContext(ctx).Where(
		u.query.User.ID.Eq(userID),
	).First()
}

func (u *UserDao) UpdateAvatar(ctx context.Context, userID int64, iconURI string) error {
	_, err := u.query.User.WithContext(ctx).Where(
		u.query.User.ID.Eq(userID),
	).Updates(map[string]interface{}{
		"icon_uri":   iconURI,
		"updated_at": time.Now().UnixMilli(),
	})
	return err
}

func (u *UserDao) CheckUniqueNameExist(ctx context.Context, name string) (bool, error) {
	_, err := u.query.User.WithContext(ctx).Select(u.query.User.ID).Where(
		u.query.User.Name.Eq(name),
	).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateUser Create a new user
func (u *UserDao) CreateUser(ctx context.Context, user *model.User) error {
	return u.query.User.WithContext(ctx).Create(user)
}
