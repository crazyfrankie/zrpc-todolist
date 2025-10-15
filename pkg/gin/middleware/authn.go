package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/crazyfrankie/zrpc/metadata"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/response"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/conv"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
)

type AuthnHandler struct {
	noAuthPaths map[string]struct{}
	authClient  auth.AuthServiceClient
}

func NewAuthnHandler(authClient auth.AuthServiceClient) (*AuthnHandler, error) {
	return &AuthnHandler{authClient: authClient, noAuthPaths: make(map[string]struct{})}, nil
}

func (h *AuthnHandler) IgnorePath(paths []string) *AuthnHandler {
	for _, path := range paths {
		h.noAuthPaths[path] = struct{}{}
	}
	return h
}

func (h *AuthnHandler) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		md := metadata.New(map[string]string{
			"user_agent": c.Request.UserAgent(),
		})

		currentPath := c.Request.URL.Path

		if _, ok := h.noAuthPaths[currentPath]; ok {
			c.Request = c.Request.WithContext(h.storeUserInfo(c, md))
			c.Next()
			return
		}

		accessToken, err := getAccessToken(c)
		if err != nil {
			response.Unauthorized(c)
			return
		}
		parseRes, err := h.authClient.ParseToken(c.Request.Context(), &auth.ParseTokenRequest{Token: accessToken})
		if err == nil {
			md.Append("user_id", conv.Int64ToStr(parseRes.GetUserID()))
			c.Request = c.Request.WithContext(h.storeUserInfo(c, md))

			c.Next()
			return
		}

		response.Unauthorized(c)
	}
}

func (h *AuthnHandler) storeUserInfo(ctx context.Context, md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}

func getAccessToken(c *gin.Context) (string, error) {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		return "", errors.New("no auth")
	}

	strs := strings.Split(tokenHeader, " ")
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", errors.New("header is invalid")
	}

	if strs[1] == "" {
		return "", errors.New("jwt is empty")
	}

	return strs[1], nil
}
