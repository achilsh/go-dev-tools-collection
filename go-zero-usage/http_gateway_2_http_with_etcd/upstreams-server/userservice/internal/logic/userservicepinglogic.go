package logic

import (
	"context"

	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserservicePingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserservicePingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserservicePingLogic {
	return &UserservicePingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserservicePingLogic) UserservicePing(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
