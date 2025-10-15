package token

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

type Token interface {
	GenerateToken(uid int64) ([]string, error)
	ParseToken(token string) (*Claims, error)
	TryRefresh(refresh string) ([]string, int64, error)
	CleanToken(ctx context.Context, uid int64) error
}
