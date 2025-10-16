package service

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/token"
)

type Components struct {
	TokenGen token.Token
}

type authImpl struct {
	*Components
}

func NewAuthDomain(c *Components) Auth {
	return &authImpl{c}
}

func (a *authImpl) GenerateToken(ctx context.Context, uid int64) ([]string, error) {
	tokens, err := a.TokenGen.GenerateToken(uid)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (a *authImpl) ParseToken(ctx context.Context, token string) (*token.Claims, error) {
	claims, err := a.TokenGen.ParseToken(token)
	if err != nil {
		return nil, err
	}

	return claims, err
}

func (a *authImpl) RefreshToken(ctx context.Context, refreshToken string) ([]string, int64, error) {
	tokens, userID, err := a.TokenGen.TryRefresh(refreshToken)
	if err != nil {
		return nil, 0, err
	}

	return tokens, userID, nil
}

func (a *authImpl) CleanToken(ctx context.Context, userID int64) error {
	return a.TokenGen.CleanToken(ctx, userID)
}
