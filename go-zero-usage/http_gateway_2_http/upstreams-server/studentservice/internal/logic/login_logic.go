package logic

import (
	"context"

	"http_gateway_demo/upstreams-server/studentservice/internal/svc"
	"http_gateway_demo/upstreams-server/studentservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("login req: %+v", *req)
	resp = &types.LoginResponse{
		Code:    123,
		Message: "this is message",
	}
	return
}
