package token

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/cache"
	"github.com/crazyfrankie/zrpc-todolist/infra/contract/token"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

const (
	RefreshPrefix = "refresh_token"
)

type TokenService struct {
	cmd       cache.Cmdable
	signAlgo  string
	secretKey *rsa.PrivateKey
	publicKey *rsa.PublicKey
}

func New(cmd cache.Cmdable) (token.Token, error) {
	signAlgo := os.Getenv(consts.JWTSignAlgo)
	secretPath := os.Getenv(consts.JWTSecretKey)
	publicPath := os.Getenv(consts.JWTPublicKey)

	privateKey, err := os.ReadFile(secretPath)
	if err != nil {
		return nil, err
	}
	private, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	publicKey, _ := os.ReadFile(publicPath)
	public, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	return &TokenService{cmd: cmd, signAlgo: signAlgo, secretKey: private, publicKey: public}, nil
}

func (s *TokenService) GenerateToken(uid int64) ([]string, error) {
	res := make([]string, 2)
	access, err := s.newToken(uid, time.Hour*2)
	if err != nil {
		return res, err
	}
	res[0] = access
	refresh, err := s.newToken(uid, time.Hour*24*30)
	if err != nil {
		return res, err
	}
	res[1] = refresh

	// set refresh in redis
	key := refreshKey(uid)

	err = s.cmd.Set(context.Background(), key, refresh, time.Hour*24*30).Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *TokenService) newToken(uid int64, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &token.Claims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}
	tk := jwt.NewWithClaims(jwt.GetSigningMethod(s.signAlgo), claims)
	str, err := tk.SignedString(s.secretKey)

	return str, err
}

func (s *TokenService) ParseToken(tk string) (*token.Claims, error) {
	t, err := jwt.ParseWithClaims(tk, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*token.Claims)
	if !ok {
		return nil, errors.New("jwt is invalid")
	}

	return claims, nil
}

func (s *TokenService) TryRefresh(refresh string) ([]string, int64, error) {
	refreshClaims, err := s.ParseToken(refresh)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid refresh jwt")
	}

	res, err := s.cmd.Get(context.Background(), refreshKey(refreshClaims.UID)).Result()
	if err != nil || res != refresh {
		return nil, 0, errors.New("jwt invalid or revoked")
	}

	access, err := s.newToken(refreshClaims.UID, time.Hour*2)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	issat, _ := refreshClaims.GetIssuedAt()
	expire, _ := refreshClaims.GetExpirationTime()
	if expire.Sub(now) < expire.Sub(issat.Time)/3 {
		// try refresh
		refresh, err = s.newToken(refreshClaims.UID, time.Hour*24*30)
		err = s.cmd.Set(context.Background(), refreshKey(refreshClaims.UID), refresh, time.Hour*24*30).Err()
		if err != nil {
			return nil, 0, err
		}
	}

	return []string{access, refresh}, refreshClaims.UID, nil
}

func (s *TokenService) CleanToken(ctx context.Context, uid int64) error {
	return s.cmd.Del(ctx, refreshKey(uid)).Err()
}

func refreshKey(uid int64) string {
	return fmt.Sprintf("%s:%d", RefreshPrefix, uid)
}
