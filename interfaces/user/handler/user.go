package handler

import (
	"io"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/interfaces/user/model"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/response"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/conv"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
	"github.com/crazyfrankie/zrpc-todolist/protocol/user"
)

type UserHandler struct {
	userClient user.UserServiceClient
}

func NewUserHandler(userClient user.UserServiceClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

func (h *UserHandler) RegisterRoute(r *gin.RouterGroup) {
	userGroup := r.Group("user")
	{
		userGroup.POST("register", h.Register())
		userGroup.POST("login", h.Login())
		userGroup.GET("logout", h.Logout())
		userGroup.GET("profile", h.GetUserInfo())
		userGroup.POST("avatar", h.UpdateAvatar())
		userGroup.POST("reset-password", h.ResetPassword())
	}
}

func (h *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.UserRegisterReq
		if err := c.ShouldBind(&req); err != nil {
			response.InvalidParamError(c, err.Error())
			return
		}

		res, err := h.userClient.Register(c.Request.Context(), &user.RegisterRequest{
			Name:     req.Name,
			Password: req.Password,
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.SetAuthorization(c, res.Data.AccessToken, res.Data.RefreshToken)

		response.Success(c, userDTO2VO(res.Data))
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.UserLoginReq
		if err := c.ShouldBind(&req); err != nil {
			response.InvalidParamError(c, err.Error())
			return
		}

		res, err := h.userClient.Login(c.Request.Context(), &user.LoginRequest{
			Name:     req.Name,
			Password: req.Password,
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.SetAuthorization(c, res.Data.AccessToken, res.Data.RefreshToken)

		response.Success(c, userDTO2VO(res.Data))
	}
}

func (h *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := h.userClient.Logout(c.Request.Context(), &user.LogoutRequest{})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, nil)
	}
}

func (h *UserHandler) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.userClient.GetUserInfo(c.Request.Context(), &user.GetUserInfoRequest{})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, userDTO2VO(res.Data))
	}
}

func (h *UserHandler) UpdateAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("avatar")
		if err != nil {
			logs.CtxErrorf(c.Request.Context(), "Get Avatar Fail failed, err=%v", err)
			response.InvalidParamError(c, "missing avatar file")
			return
		}

		// Check file type
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			response.InvalidParamError(c, "invalid file type, only image allowed")
			return
		}

		// Read file content
		src, err := file.Open()
		if err != nil {
			response.InternalServerError(c, err)
			return
		}
		defer src.Close()

		fileContent, err := io.ReadAll(src)
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		res, err := h.userClient.UpdateAvatar(c.Request.Context(), &user.UpdateAvatarRequest{
			Avatar:   fileContent,
			MimeType: file.Header.Get("Content-Type"),
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, res)
	}
}

func (h *UserHandler) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.UserResetPasswordReq
		if err := c.ShouldBind(&req); err != nil {
			response.InvalidParamError(c, err.Error())
			return
		}

		_, err := h.userClient.ResetPassword(c.Request.Context(), &user.ResetPasswordRequest{
			Name:     req.Name,
			Password: req.Password,
		})
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		response.Success(c, nil)
	}
}

func userDTO2VO(userDto *user.User) *model.UserInfoResp {
	return &model.UserInfoResp{
		UserID:         conv.Int64ToStr(userDto.UserID),
		Name:           userDto.Name,
		Avatar:         userDto.AvatarUrl,
		UserCreateTime: userDto.UserCreateTime,
	}
}
