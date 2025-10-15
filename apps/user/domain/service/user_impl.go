package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/entity"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/internal/dal/model"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/repository"
	"github.com/crazyfrankie/zrpc-todolist/infra/contract/idgen"
	"github.com/crazyfrankie/zrpc-todolist/infra/contract/storage"
	"github.com/crazyfrankie/zrpc-todolist/pkg/errorx"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/conv"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
	"github.com/crazyfrankie/zrpc-todolist/types/errno"
)

type Components struct {
	UserRepo repository.UserRepository
	IconOSS  storage.Storage
	IDGen    idgen.IDGenerator
}

type userImpl struct {
	*Components
}

func NewUserDomain(c *Components) User {
	return &userImpl{c}
}

func (u *userImpl) Create(ctx context.Context, req *CreateUserRequest) (*entity.User, error) {
	if req.Name != "" {
		exist, err := u.UserRepo.CheckUniqueNameExist(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errorx.New(errno.ErrUserUniqueNameAlreadyExistCode, errorx.KV("name", req.Name))
		}
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userID, err := u.IDGen.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	newUser := &model.User{
		ID:       userID,
		IconURI:  consts.UserIconURI,
		Name:     req.Name,
		Password: hashedPassword,
	}

	err = u.UserRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("insert user failed: %w", err)
	}

	iconURL, err := u.IconOSS.GetObjectUrl(ctx, newUser.IconURI)
	if err != nil {
		return nil, fmt.Errorf("get icon url failed: %w", err)
	}

	return userPO2DO(newUser, iconURL), nil
}

func (u *userImpl) Login(ctx context.Context, name, password string) (*entity.User, error) {
	userModel, exist, err := u.UserRepo.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errorx.New(errno.ErrUserInfoInvalidCode)
	}

	valid := verifyPassword(password, userModel.Password)
	if !valid {
		return nil, errorx.New(errno.ErrUserInfoInvalidCode)
	}

	resURL, err := u.IconOSS.GetObjectUrl(ctx, userModel.IconURI)
	if err != nil {
		return nil, err
	}

	return userPO2DO(userModel, resURL), nil
}

func (u *userImpl) ResetPassword(ctx context.Context, name, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	return u.UserRepo.UpdatePassword(ctx, name, hashedPassword)
}

func (u *userImpl) GetUserInfo(ctx context.Context, userID int64) (user *entity.User, err error) {
	if userID <= 0 {
		return nil, errorx.New(errno.ErrUserInvalidParamCode,
			errorx.KVf("msg", "invalid user id : %d", userID))
	}

	userModel, err := u.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resURL, err := u.IconOSS.GetObjectUrl(ctx, userModel.IconURI)
	if err != nil {
		return nil, err
	}

	return userPO2DO(userModel, resURL), nil
}

func (u *userImpl) UpdateAvatar(ctx context.Context, userID int64, ext string, imagePayload []byte) (url string, err error) {
	avatarKey := "user_avatar/" + conv.Int64ToStr(userID) + "." + ext
	err = u.IconOSS.PutObject(ctx, avatarKey, imagePayload)
	if err != nil {
		return "", err
	}

	err = u.UserRepo.UpdateAvatar(ctx, userID, avatarKey)
	if err != nil {
		return "", err
	}

	url, err = u.IconOSS.GetObjectUrl(ctx, avatarKey)
	if err != nil {
		return "", err
	}

	return url, nil
}

func userPO2DO(model *model.User, iconURL string) *entity.User {
	return &entity.User{
		UserID:    model.ID,
		Name:      model.Name,
		IconURI:   model.IconURI,
		IconURL:   iconURL,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func hashPassword(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPass), nil
}

func verifyPassword(password, encodedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
