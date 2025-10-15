package application

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/apps/auth/domain/service"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
)

type AuthApplicationService struct {
	authDomain service.Auth
	auth.UnimplementedAuthServiceServer
}

func NewAuthApplicationService(authDomain service.Auth) *AuthApplicationService {
	return &AuthApplicationService{authDomain: authDomain}
}

func (a *AuthApplicationService) GenerateToken(ctx context.Context, req *auth.GenerateTokenRequest) (*auth.GenerateTokenResponse, error) {
	tokens, err := a.authDomain.GenerateToken(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}

	return &auth.GenerateTokenResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	}, nil
}

func (a *AuthApplicationService) ParseToken(ctx context.Context, req *auth.ParseTokenRequest) (*auth.ParseTokenResponse, error) {
	claims, err := a.authDomain.ParseToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &auth.ParseTokenResponse{UserID: claims.UID}, nil
}

func (a *AuthApplicationService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	tokens, userID, err := a.authDomain.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &auth.RefreshTokenResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
		UserID:       userID,
	}, nil
}
