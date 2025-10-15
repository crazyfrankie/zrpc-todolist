package service

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/token"
)

type Auth interface {
	GenerateToken(ctx context.Context, uid int64) ([]string, error)
	ParseToken(ctx context.Context, token string) (*token.Claims, error)
	RefreshToken(ctx context.Context, refreshToken string) ([]string, int64, error)
}
