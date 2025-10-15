package application

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/entity"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/service"
	"github.com/crazyfrankie/zrpc-todolist/pkg/errorx"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/ctxutil"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/protocol/user"
	"github.com/crazyfrankie/zrpc-todolist/types/errno"
)

type UserApplicationService struct {
	userDomain service.User
	authClient auth.AuthServiceClient

	user.UnimplementedUserServiceServer
}

func NewUserApplicationService(userDomain service.User, authClient auth.AuthServiceClient) *UserApplicationService {
	return &UserApplicationService{userDomain: userDomain, authClient: authClient}
}

func (u *UserApplicationService) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	userInfo, err := u.userDomain.Create(ctx, &service.CreateUserRequest{
		Password: req.GetPassword(),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	userInfo, err = u.userDomain.Login(ctx, req.GetName(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	tkRes, err := u.authClient.GenerateToken(ctx, &auth.GenerateTokenRequest{
		UserID: userInfo.UserID,
	})
	if err != nil {
		return nil, err
	}

	data := userDO2DTO(userInfo)
	data.AccessToken = tkRes.AccessToken
	data.RefreshToken = tkRes.RefreshToken

	metrics.UserRegisterCounter.Add(1)

	return &user.RegisterResponse{
		Data: data,
	}, nil
}

func (u *UserApplicationService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	userInfo, err := u.userDomain.Login(ctx, req.GetName(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	tkRes, err := u.authClient.GenerateToken(ctx, &auth.GenerateTokenRequest{
		UserID: userInfo.UserID,
	})
	if err != nil {
		return nil, err
	}

	data := userDO2DTO(userInfo)
	data.AccessToken = tkRes.AccessToken
	data.RefreshToken = tkRes.RefreshToken

	metrics.UserLoginCounter.Add(1)

	return &user.LoginResponse{
		Data: data,
	}, nil
}

func (u *UserApplicationService) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	userID := ctxutil.MustGetUserIDFromCtx(ctx)

	userInfo, err := u.userDomain.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &user.GetUserInfoResponse{Data: userDO2DTO(userInfo)}, nil
}

func (u *UserApplicationService) UpdateAvatar(ctx context.Context, req *user.UpdateAvatarRequest) (*user.UpdateAvatarResponse, error) {
	var ext string
	var err error
	switch req.GetMimeType() {
	case "image/jpeg", "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "image/webp":
		ext = "webp"
	default:
		return nil, errorx.WrapByCode(err, errno.ErrUserInvalidParamCode,
			errorx.KV("msg", "unsupported image type"))
	}

	userID := ctxutil.MustGetUserIDFromCtx(ctx)

	iconUrl, err := u.userDomain.UpdateAvatar(ctx, userID, ext, req.GetAvatar())
	if err != nil {
		return nil, err
	}

	return &user.UpdateAvatarResponse{AvatarUrl: iconUrl}, nil
}

func (u *UserApplicationService) ResetPassword(ctx context.Context, req *user.ResetPasswordRequest) (*user.ResetPasswordResponse, error) {
	err := u.userDomain.ResetPassword(ctx, req.GetName(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &user.ResetPasswordResponse{}, nil
}

func (u *UserApplicationService) Logout(ctx context.Context, req *user.LogoutRequest) (*user.LogoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func userDO2DTO(userDo *entity.User) *user.User {
	return &user.User{
		UserID:    userDo.UserID,
		Name:      userDo.Name,
		AvatarUrl: userDo.IconURL,

		UserCreateTime: userDo.CreatedAt / 1000,
	}
}
