package logic

import (
	"context"

	"gen_jwt_token_and_check/internal/svc"
	"gen_jwt_token_and_check/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	// todo: add your logic here and delete this line
	tokenInfo, ok := l.ctx.Value("user_id").(string)
	if ok {
		l.Logger.Debugf("token info: %v", tokenInfo)
	}
	return
}
