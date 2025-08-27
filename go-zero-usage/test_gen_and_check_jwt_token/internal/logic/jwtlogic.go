package logic

import (
	"context"
	"time"

	"gen_jwt_token_and_check/internal/svc"
	"gen_jwt_token_and_check/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type JwtLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJwtLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JwtLogic {
	return &JwtLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JwtLogic) Jwt(req *types.JwtTokenRequest) (resp *types.JwtTokenResponse, err error) {
	// todo: add your logic here and delete this line

	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire

	now := time.Now().Unix()
	accessToken, err := l.GenToken(now,
		l.svcCtx.Config.JwtAuth.AccessSecret,
		map[string]interface{}{
			"user_id": req.UserId},
		accessExpire)
	if err != nil {
		return nil, err
	}

	return &types.JwtTokenResponse{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil

}

func (l *JwtLogic) GenToken(
	iat int64,
	secretKey string,
	payloads map[string]interface{},
	seconds int64,
) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
